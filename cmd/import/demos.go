/*
Copyright © 2023 Alexander Orban <alexander.orban@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package _import

import (
	"encoding/json"
	"errors"
	"fmt"
	icm_orm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"snowlastic-cli/pkg/es"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demos",
	Short: "Index a pre-defined list of documents into the `demo` index",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err        error
			caCert     []byte
			caCertPath string
			cfg        es.ElasticClientConfig

			demos []Demo
			b     []byte
			docs  = make(chan icm_orm.ICMEntity, es.BulkInsertSize)
			c     *elasticsearch.Client

			indexName = "demos"

			numErrors  int64
			numIndexed int64
		)

		// generate the CA Certificate bytes needed for the elasticsearch Config
		caCertPath = viper.GetString("elasticCaCertPath")
		caCert, err = os.ReadFile(caCertPath)
		if err != nil {
			return err
		}
		cfg = es.ElasticClientConfig{
			Addresses: []string{fmt.Sprintf(
				"https://%s:%s",
				viper.GetString("elasticUrl"),
				viper.GetString("elasticPort"),
			)},
			User:         viper.GetString("elasticUser"),
			Pass:         viper.GetString("elasticPassword"),
			ApiKey:       viper.GetString("elasticApiKey"),
			ServiceToken: viper.GetString("elasticServiceToken"),
			CaCert:       caCert,
		}
		// Generate the client
		c, err = es.NewElasticClient(&cfg)
		if err != nil {
			return err
		}

		// Get demos array
		log.Println("reading demos json")
		b = []byte(_demos)
		log.Println("unmarshalling demos into []Demo")
		err = json.Unmarshal(b, &demos)
		if err != nil {
			return errors.New(fmt.Sprintf("error in unmarshalling demos json: %s", err))
		}

		start := time.Now().UTC()
		go func() {
			// we cannot use _, demo := range demos here, since we need to pass
			// a pointer to the element as an ICMEntity. When using _, demo := range demos
			// the pointer will always point to the last document in the list.
			// Instead we point directly to the entry in the slice of demos we've created above
			for i := range demos {
				docs <- icm_orm.ICMEntity(&demos[i])
			}
			close(docs)
		}()

		batches := es.BatchEntities(docs, es.BulkInsertSize)
		numIndexed, numErrors, err = es.BulkImport(c, batches, indexName, int64(math.Ceil(float64(len(demos))/es.BulkInsertSize)))
		if err != nil {
			return err
		}

		dur := time.Since(start)
		if numErrors > 0 {
			return errors.New(fmt.Sprintf(
				"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				humanize.Comma(int64(numErrors)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			))
		} else {
			log.Printf(
				"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			)
		}
		return nil
	},
}

func init() {
	importCmd.AddCommand(demoCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// demoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// demoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Demo struct {
	ID          int    `json:"id"`
	SearchTerm  string `json:"search-term"`
	Value       string `json:"value"`
	ShouldMatch bool   `json:"should-match"`
}

func (d *Demo) IsICMEntity() bool { return true }
func (d *Demo) GetID() string     { return strconv.Itoa(d.ID) }

const _demos string = `[
  {
    "id": 1,
    "search-term": "against /5 law!",
    "value": "this might actually be against the law",
    "should-match": true
  },
  {
    "id": 2,
    "search-term": "against /5 law!",
    "value": "Law firms are against fraud",
    "should-match": true
  },
  {
    "id": 3,
    "search-term": "against /5 law!",
    "value": "this might be against several laws",
    "should-match": true
  },
  {
    "id": 4,
    "search-term": "against /5 law!",
    "value": "against isn't a term contained in the law",
    "should-match": false
  },
  {
    "id": 5,
    "search-term": "agent or agency",
    "value": "noone has survived an encounter with an agent. Noone",
    "should-match": true
  },
  {
    "id": 6,
    "search-term": "agent or agency",
    "value": "the national security agency is located in Maryland",
    "should-match": true
  },
  {
    "id": 7,
    "search-term": "agent or agency",
    "value": "One should always strive to maintain agency over oneself",
    "should-match": true
  },
  {
    "id": 8,
    "search-term": "agent or agency",
    "value": "Be a postive agent for change",
    "should-match": true
  },
  {
    "id": 9,
    "search-term": "\"appear legitimate\"",
    "value": "we need this to appear legitimate",
    "should-match": true
  },
  {
    "id": 10,
    "search-term": "\"appear legitimate\"",
    "value": "This doesn't appear to be legitimate",
    "should-match": false
  },
  {
    "id": 11,
    "search-term": "audit",
    "value": "We cannot let internal audit get wind of this",
    "should-match": true
  },
  {
    "id": 12,
    "search-term": "audit",
    "value": "A good audit would reveal the problem",
    "should-match": true
  },
  {
    "id": 13,
    "search-term": "audit",
    "value": "No real auditing procedure takes place on these",
    "should-match": false
  },
  {
    "id": 14,
    "search-term": "black /1 deal!",
    "value": "Russian black market dealers",
    "should-match": true
  },
  {
    "id": 15,
    "search-term": "black /1 deal!",
    "value": "A \"black book deal\" or an off-book account",
    "should-match": true
  },
  {
    "id": 16,
    "search-term": "black /1 deal!",
    "value": "Poker players looking for dealings of black aces, i.e clubs or spades",
    "should-match": true
  },
  {
    "id": 17,
    "search-term": "bogus /1 company!",
    "value": "A fake company is sometimes called a \"bogus company\"",
    "should-match": true
  },
  {
    "id": 18,
    "search-term": "bogus /1 company!",
    "value": "all the invoices from these companies are bogus",
    "should-match": true
  },
  {
    "id": 19,
    "search-term": "bonus!",
    "value": "I hope we get paid our bonuses this year",
    "should-match": true
  },
  {
    "id": 20,
    "search-term": "bonus!",
    "value": "A bonus for exceptional performance is sometimes appropriate",
    "should-match": true
  },
  {
    "id": 21,
    "search-term": "bonus!",
    "value": "We could class our payments to the vendor ass a bonus?",
    "should-match": true
  },
  {
    "id": 22,
    "search-term": "break! /5 law!",
    "value": "It really feels like we might be breaking the law",
    "should-match": true
  },
  {
    "id": 23,
    "search-term": "break! /5 law!",
    "value": "It wouldn't be good to break any more laws",
    "should-match": true
  },
  {
    "id": 24,
    "search-term": "break! /5 law!",
    "value": "Hoe many laws could possibly be broken?",
    "should-match": false
  },
  {
    "id": 25,
    "search-term": "brib!",
    "value": "could this be a bribe?",
    "should-match": true
  },
  {
    "id": 26,
    "search-term": "brib!",
    "value": "Some countries require bribes are a part of doing business!",
    "should-match": true
  },
  {
    "id": 27,
    "search-term": "brib!",
    "value": "Bribing SOE officials is certainly a problem",
    "should-match": true
  },
  {
    "id": 28,
    "search-term": "broke! /5 law!",
    "value": "Hoe many laws could possibly be broken?",
    "should-match": true
  },
  {
    "id": 29,
    "search-term": "broke! /5 law!",
    "value": "I think he may have broke the law with this one",
    "should-match": true
  },
  {
    "id": 30,
    "search-term": "broke! /5 law!",
    "value": "law within a few more words of extra long distiance from broken",
    "should-match": false
  },
  {
    "id": 31,
    "search-term": "buyoff! or (buy! /1 off) or (buy /1 off!)",
    "value": "This seems like a good time to buy their offerings",
    "should-match": true
  },
  {
    "id": 32,
    "search-term": "buyoff! or (buy! /1 off) or (buy /1 off!)",
    "value": "We should considering buying something off the shelf, rather than building it",
    "should-match": true
  },
  {
    "id": 33,
    "search-term": "buyoff! or (buy! /1 off) or (buy /1 off!)",
    "value": "setting up good software would allow our team to buy items while offline, not just in the office",
    "should-match": false
  },
  {
    "id": 34,
    "search-term": "cash",
    "value": "We should think about cashing out from the poker table",
    "should-match": true
  },
  {
    "id": 35,
    "search-term": "cash",
    "value": "make an entry in the cashbook",
    "should-match": true
  },
  {
    "id": 36,
    "search-term": "cash",
    "value": "we are looking to hire a cashier on a part-time basis",
    "should-match": false
  },
  {
    "id": 37,
    "search-term": "charit!",
    "value": "a donation to a charity benefits our business' reputaiton",
    "should-match": true
  },
  {
    "id": 38,
    "search-term": "charit!",
    "value": "a highly uncharitable view on external sales",
    "should-match": true
  },
  {
    "id": 39,
    "search-term": "charit!",
    "value": "some charities represent a risk on funding unwanted activities",
    "should-match": true
  },
  {
    "id": 40,
    "search-term": "compliance",
    "value": "An uncomplimentary perspective on off-the-shelf products",
    "should-match": false
  },
  {
    "id": 41,
    "search-term": "compliance",
    "value": "noncompliance in controls creates hazards",
    "should-match": true
  },
  {
    "id": 42,
    "search-term": "compliance",
    "value": "heavy handed reliance on trianing may create overcompliance in some circumstances",
    "should-match": true
  },
  {
    "id": 43,
    "search-term": "corrupt!",
    "value": "incorruptable behavior is what we hope to create",
    "should-match": false
  },
  {
    "id": 44,
    "search-term": "corrupt!",
    "value": "financial criminals hired into an organization may act corruptively",
    "should-match": true
  },
  {
    "id": 45,
    "search-term": "corrupt!",
    "value": "anticorruption data scientists are a critical part of the FCPA framework",
    "should-match": true
  },
  {
    "id": 46,
    "search-term": "discount /10 increase",
    "value": "we want to provide a discount to our customer in order to increase overall sales",
    "should-match": true
  },
  {
    "id": 47,
    "search-term": "discount /10 increase",
    "value": "can we increase the discount to this customer?",
    "should-match": true
  },
  {
    "id": 48,
    "search-term": "discount /10 increase",
    "value": "no. increasing the discount to this customer is unliely to generate more sales",
    "should-match": false
  },
  {
    "id": 49,
    "search-term": "donat!",
    "value": "please submit an approval request for a donation to the ABB charity",
    "should-match": true
  },
  {
    "id": 50,
    "search-term": "donat!",
    "value": "The legions, uninflamed by party zeal, were allured into civil war by liberal donatives, and still more liberal promises.",
    "should-match": true
  },
  {
    "id": 51,
    "search-term": "donat!",
    "value": "Sometimes the means would not be just ready to the day, but almost invariably donators paid honorably.",
    "should-match": true
  },
  {
    "id": 52,
    "search-term": "fee",
    "value": "Unless he gets a respite from another court, the rapper will appear on just 12 ballots, three of them in states where he needed only to pay a fee for access.",
    "should-match": true
  },
  {
    "id": 53,
    "search-term": "fee",
    "value": "The firm has 45 commitments from investors, including many in the $1 billion range, it said without specifying whether the money is heading into the high-fee Pure Alpha hedge funds or low fee long-only products.",
    "should-match": true
  },
  {
    "id": 54,
    "search-term": "fee",
    "value": "Now Arm is being acquired by one of those competitors, which may use its position to hike licensing and royalty fees on its rivals or to deny them access to the latest technology.",
    "should-match": true
  },
  {
    "id": 55,
    "search-term": "intervene or intervention",
    "value": "The weekly reports have always been published by scientists and other public health professionals alone, without other branches of the government intervening.",
    "should-match": true
  },
  {
    "id": 56,
    "search-term": "intervene or intervention",
    "value": "Hancock said the new rule will be “rigorously enforced by police,” who currently have no powers to intervene when up to 30 people gather.",
    "should-match": true
  },
  {
    "id": 57,
    "search-term": "intervene or intervention",
    "value": "With so much research to wade through, it’s hard to know what to trust — and I say that as someone who makes a living researching what types of interventions motivate people to change their behaviors.",
    "should-match": true
  },
  {
    "id": 58,
    "search-term": "intervene or intervention",
    "value": "Researchers have rightly realized that individual variation is just as important as the average response to an intervention.",
    "should-match": true
  },
  {
    "id": 59,
    "search-term": "money",
    "value": "Combining that with advertisers’ increased upfront cancelation options, the money committed to traditional TV could wind up going to streaming.",
    "should-match": true
  },
  {
    "id": 60,
    "search-term": "money",
    "value": "Noonan said his daughters sometimes work from an office in the Phoenix area and are classified as independent contractors, not earning “horrible money” but also not making minimum wage.",
    "should-match": true
  },
  {
    "id": 61,
    "search-term": "money",
    "value": "“Profits” were returned to early investors with monies from newer victims.",
    "should-match": true
  },
  {
    "id": 62,
    "search-term": "red /1 flag!",
    "value": "The celebrities touting the app to their millions of followers have sent up red flags to health experts.",
    "should-match": true
  },
  {
    "id": 63,
    "search-term": "red /1 flag!",
    "value": "A red flag meant that there had been no food at the protest center.",
    "should-match": true
  },
  {
    "id": 64,
    "search-term": "red /1 flag!",
    "value": "That day, a prominent march was planned at the historic Red Fort, where India’s prime minister traditionally hoists the flag on Independence Day.",
    "should-match": false
  },
  {
    "id": 65,
    "search-term": "\"appear legitimate\"",
    "value": "The story of fluoridation reads like a postmodern fable, and the moral is clear: a scientific discovery might seem like a legitimate boon.",
    "should-match": false
  },
  {
    "id": 66,
    "search-term": "\"appear legitimate\"",
    "value": "Again, the difference can seem legitimatly subtle and sound more like splitting hairs, but the difference is important.",
    "should-match": true
  },
  {
    "id": 67,
    "search-term": "\"appear legitimate\"",
    "value": "So far, all the players seemed legitimate.",
    "should-match": true
  },
  {
    "id": 68,
    "search-term": "\"special /1 commission!\"",
    "value": "The special commission issued a request for proposals but never awarded a contract.",
    "should-match": true
  },
  {
    "id": 69,
    "search-term": "\"special /1 commission!\"",
    "value": "The commission has yet to formally award the special work to any contractors.",
    "should-match": false
  },
  {
    "id": 70,
    "search-term": "\"special /1 commission!\"",
    "value": "A special commission will be awarded if the work is completed within 24 hours",
    "should-match": true
  },
  {
    "id": 71,
    "search-term": "Advisor!",
    "value": "Calvert, meanwhile, did not disclose that he once held an unpaid position on the advisory board of Hayungs’ company, MRG Medical.",
    "should-match": true
  },
  {
    "id": 72,
    "search-term": "Advisor!",
    "value": "The report also recommends establishing an international scientific advisory board to evaluate the state of the technology and consult on applications to do such heritable or germline editing.",
    "should-match": true
  },
  {
    "id": 73,
    "search-term": "Advisor!",
    "value": "Directing duties are split among Arena’s formidable out artistic director Molly Smith, deputy artistic director Seema Sueko, and senior artistic adviser Anita Maynard-Losh, along with local directors Paige Hernandez and director Psalmayene 24.",
    "should-match": true
  },
  {
    "id": 74,
    "search-term": "Black account!",
    "value": "we've hired more than 300 african-american (black) accountants to satisfy our ESG initiatives",
    "should-match": true
  },
  {
    "id": 75,
    "search-term": "Black account!",
    "value": "make an entry into the black account books",
    "should-match": true
  },
  {
    "id": 76,
    "search-term": "Carrot!",
    "value": "if the carrot doesn't work, then use the stick",
    "should-match": true
  },
  {
    "id": 77,
    "search-term": "Carrot!",
    "value": "He had a mop of carroty hair, and on top of it was a little plaid cap that looked as though it was lost in the wilderness.",
    "should-match": true
  },
  {
    "id": 78,
    "search-term": "Carrot!",
    "value": "Over behind the stove was a tall, awkward boy with carroty hair and small, dark eyes set much aslant in the saddest of faces.",
    "should-match": true
  },
  {
    "id": 79,
    "search-term": "Cash!",
    "value": "If you put all the green blocks in a lineup and make me guess which is asparagus square and which is pistachio and cashew square, I’ll have no idea.",
    "should-match": true
  },
  {
    "id": 80,
    "search-term": "Cash!",
    "value": "Fresh pork stewed with parsnips; turnips; winter-squash or cashaw—Apple dumplings.",
    "should-match": true
  },
  {
    "id": 81,
    "search-term": "Cash!",
    "value": "In many of the cities and large towns, some credit grocers have adopted what is called the \"cash-and-carry plan.\"",
    "should-match": true
  },
  {
    "id": 82,
    "search-term": "Circumven!",
    "value": "This is a way to circumvent browser blockades to ensure that ad platforms can access the data advertisers rely on to measure campaign performance now and in the future when a website’s server can no longer set trackers on the user’s browser.",
    "should-match": true
  },
  {
    "id": 83,
    "search-term": "Circumven!",
    "value": "The idea of giving freely of what you no longer need rather than tossing it—that is, the very essence of a hiker box—circumvents the self-replicating loop of infinite consumption and waste.",
    "should-match": true
  },
  {
    "id": 84,
    "search-term": "Circumven!",
    "value": "Yet in reality, those protections turned out to be either ineffective or easy to circumvent, our stories showed.",
    "should-match": true
  },
  {
    "id": 85,
    "search-term": "Delete!",
    "value": "Still, negative association in the card environment had a deleterious effect, said Manatt.",
    "should-match": true
  },
  {
    "id": 86,
    "search-term": "Delete!",
    "value": "One 2018 analysis suggests that around 1 out of every 100 people has deleterious mosaic genetic difference that affects “sizable brain regions.”",
    "should-match": true
  },
  {
    "id": 87,
    "search-term": "Delete!",
    "value": "So it has also ordered the police to take steps to ensure Clearview deletes the data.",
    "should-match": true
  },
  {
    "id": 88,
    "search-term": "Delete!",
    "value": "Often material will remain in the company’s servers even if it has been deleted off their sites.",
    "should-match": true
  },
  {
    "id": 89,
    "search-term": "Delete!",
    "value": "I really need to delete some of my old tweets",
    "should-match": true
  },
  {
    "id": 90,
    "search-term": "Jail",
    "value": "Maybe only a federal effort to establish standards and regulate compliance to them would be necessary before we no longer have a Robert Williams, a member of any minority group, or any citizen unjustly experience a night in jail or worse.",
    "should-match": true
  },
  {
    "id": 91,
    "search-term": "Jail",
    "value": "Anyone who violated the law would be subject to a fine of $100-$500 – the equivalent of $1,700-$8,500 today – or a jail term of up to 150 days.",
    "should-match": true
  },
  {
    "id": 92,
    "search-term": "Jail",
    "value": "If a person is arrested by police, they are usually jailed until they are taken to trial. Depending on the judge’s ruling, they may be jailed again as punishment for a crime.",
    "should-match": true
  },
  {
    "id": 93,
    "search-term": "Off shore",
    "value": "our off-shore bank accounts hold the bulk of our taxable cash",
    "should-match": true
  },
  {
    "id": 94,
    "search-term": "Off shore",
    "value": "At the same time, the quality of these offshore teams is worse so its one of those things that happens over time.",
    "should-match": true
  },
  {
    "id": 95,
    "search-term": "Off shore",
    "value": "In the American modeling system’s set of simulations, the majority of lows are centered well off shore the southern Delmarva Peninsula.",
    "should-match": true
  },
  {
    "id": 96,
    "search-term": "Special! w/2 designat!",
    "value": "At the same time as Google is lobbying for limits on any special gatekeeper designations, the tech giant wants to see certain types of rules applied universally to all players.",
    "should-match": true
  },
  {
    "id": 97,
    "search-term": "Special! w/2 designat!",
    "value": "The team is working closely with the FDA and was granted a breakthrough device designation in July, which could pave the way for a human trial for treating people with paraplegia and tetraplegia.",
    "should-match": true
  },
  {
    "id": 98,
    "search-term": "Special! w/2 designat!",
    "value": "Some scientists have argued that Neowise deserves a “great special comet” designation for its brightness.",
    "should-match": true
  },
  {
    "id": 99,
    "search-term": "State-own!",
    "value": "doing business with a state-owned-entity (SOE) creates a higher risk",
    "should-match": true
  },
  {
    "id": 100,
    "search-term": "State-own!",
    "value": "oil and gas businesses are notoriously stated-owned operations",
    "should-match": true
  },
  {
    "id": 101,
    "search-term": "State-own!",
    "value": "A thumbnail version of his argument is that the initial stages of the economic transition usually involve the privatization of state-owned enterprises. ",
    "should-match": true
  },
  {
    "id": 102,
    "search-term": "VIP!",
    "value": "VIP is an informal way to refer to someone who is notable in some way and is given special treatment in a particular setting.",
    "should-match": true
  },
  {
    "id": 103,
    "search-term": "VIP!",
    "value": "Excuse me, sir. Seeing as how the V.P. is such a V.I.P., shouldn't we keep the P.C. on the Q.T.? 'Cause if it leaks to the V.C. he could end up M.I.A., and then we'd all be put on K.P. ",
    "should-match": true
  },
  {
    "id": 104,
    "search-term": "Sub-contract!",
    "value": "Auditors who examined the payments said that Intralot claimed credit for paying one local subcontractor $280,000 and another $179,000 for work that the subcontractors actually did not perform themselves, in violation of city rules.",
    "should-match": true
  },
  {
    "id": 105,
    "search-term": "Sub-contract!",
    "value": "As one example, the union supported Shannon Wait, a data-center technician employed through a sub-contractor in South Carolina, through a wrongful suspension this March for speaking to coworkers about her working conditions.",
    "should-match": true
  },
  {
    "id": 106,
    "search-term": "Sub-contract!",
    "value": "Although it employs its workers directly rather than using sub-contractors, the majority are reportedly hired on a day-to-day basis the night before via an app called “Coupunch,” or on temporary contracts that usually last a few months.",
    "should-match": true
  },
  {
    "id": 107,
    "search-term": "Supporter!",
    "value": "i've long been a supporter of technology in sports",
    "should-match": true
  },
  {
    "id": 108,
    "search-term": "Supporter!",
    "value": "The presidency and its symbols do not belong to one president or his supporters, but to all of the American people.",
    "should-match": true
  },
  {
    "id": 109,
    "search-term": "Supporter!",
    "value": "In a new statement emailed to supporters this week, she closed that door.",
    "should-match": true
  }
]
`
