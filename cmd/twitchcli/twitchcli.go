package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/leekchan/accounting"
	"github.com/manifoldco/promptui"
	"github.com/mattn/go-runewidth"
	"github.com/syphoxy/go-twitch/twitch"
)

const (
	unicodeHeavyBlackHeart            = "\u2764"
	unicodeBlackRightPointingTriangle = "\u25B6"
	unicodeHeavyCheckMark             = "\u2714"

	numberGroupSize = 3
)

type Item struct {
	Id          int
	DisplayName string
	Name        string
	Viewers     int
	Status      string
	Following   bool
}

var numberFormatter = accounting.Accounting{Symbol: "", Precision: 0}

func main() {
	game := flag.String("g", "", "get list of channels by `game` (example: 'Dota 2')")
	limit := flag.Int("l", 25, "`limit` number of results returned (0-100)")
	noprompt := flag.Bool("x", false, "do not prompt. print and exit.")
	flag.Parse()

	client := twitch.NewClient(&http.Client{})

	if *limit <= 0 {
		os.Exit(0)
	}

	if *limit > 100 {
		log.Fatal("Limit must be a value between 0 and 100")
	}

	opts := twitch.ListOptions{
		Limit: *limit,
		Game:  *game,
	}

	l, err := client.Streams.List(&opts)
	if err != nil {
		log.Fatal(err)
	}

	f, err := client.Streams.Followed(&opts)
	if err != nil {
		log.Fatal(err)
	}

	i := make([]Item, 0, len(l.Streams)+len(f.Streams))

	hoisted := make(map[int]int)
	for _, s := range f.Streams {
		i = append(i, Item{
			Id:          s.Id,
			DisplayName: s.Channel.DisplayName,
			Name:        s.Channel.Name,
			Viewers:     s.Viewers,
			Status:      strings.Replace(s.Channel.Status, "\n", " ", -1),
			Following:   true,
		})
		hoisted[s.Id] = 1
	}
	for _, s := range l.Streams {
		if _, ok := hoisted[s.Id]; !ok {
			i = append(i, Item{
				Id:          s.Id,
				DisplayName: s.Channel.DisplayName,
				Name:        s.Channel.Name,
				Viewers:     s.Viewers,
				Status:      strings.Replace(s.Channel.Status, "\n", " ", -1),
				Following:   false,
			})
		}
	}

	n, v := calcColumnWidths(i)

	renderer := renderPrompt
	if *noprompt {
		renderer = renderList
	}
	renderer(i, n, v)
}

func calcColumnWidths(i []Item) (n, v int) {
	for _, s := range i {
		if l := len(s.Name); l > n {
			n = l
		}
		l := int(math.Ceil(math.Log(float64(s.Viewers)) / math.Log(10)))
		l += int(math.Floor(float64(l) / float64(numberGroupSize)))
		if l > v {
			v = l
		}
	}
	return
}

func formatLine(i Item, n, v int) string {
	vv := numberFormatter.FormatMoneyInt(i.Viewers)
	return i.Name + strings.Repeat(" ", n-runewidth.StringWidth(i.Name)) +
		" " + strings.Repeat(" ", v-len(vv)) + vv +
		" " + strings.Replace(i.Status, "\n", "", -1)
}

func renderList(i []Item, n, v int) {
	for _, s := range i {
		l := formatLine(s, n, v)
		if s.Following {
			l = promptui.Styler(promptui.FGRed, promptui.FGBold)(unicodeHeavyBlackHeart) + " " + l
		} else {
			l = "  " + l
		}
		fmt.Println(l)
	}
}

func renderPrompt(i []Item, n, v int) {
	promptui.FuncMap["following"] = func(s Item) string {
		if s.Following {
			return unicodeHeavyBlackHeart
		}
		return " "
	}

	promptui.FuncMap["active"] = func() string {
		return unicodeBlackRightPointingTriangle
	}

	promptui.FuncMap["selected"] = func() string {
		return unicodeHeavyCheckMark
	}

	promptui.FuncMap["render"] = func(s Item) string {
		return formatLine(s, n, v)
	}

	prompt := promptui.Select{
		Label:     "Select stream",
		Items:     i,
		Size:      100,
		IsVimMode: false,
		Templates: &promptui.SelectTemplates{
			Active:   "{{ active | bold }} {{ . | following | bold | red }} {{ . | render | bold }}",
			Inactive: "  {{ . | following | bold | red }} {{ . | render }}",
			Selected: "Watching twitch.tv/{{ .DisplayName }} {{ selected | bold | green }}",
		},
	}

	x, _, err := prompt.Run()
	if err != nil {
		log.Fatal("failed to run prompt:", err)
	}

	cmd, err := exec.LookPath("streamlink")
	if err != nil {
		log.Fatal("failed to find streamlink:", err)
	}

	syscall.Exec(cmd, []string{"streamlink", "twitch.tv/" + i[x].Name, "best"}, os.Environ())
}
