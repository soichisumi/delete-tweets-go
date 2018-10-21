package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	"time"
)

// delete 20000 old tweets from given date.
func deleteTweetsBefore(client *twitter.Client, twitterID string, beforeFromDate time.Time){
	//beforeFrom := time.Date(2018, 4, 1, 0, 0, 0, 0, time.Local) // example

	lastTweetID := int64(0) // keep
	for i:= 0; i< 100; i++ {
		userTimelineParams := twitter.UserTimelineParams{
			ScreenName: twitterID, //ex: _yoyoyousei
			Count:      200,
		}

		if lastTweetID != 0 {
			userTimelineParams.MaxID = lastTweetID
		}

		tweets, _, err := client.Timelines.UserTimeline(&userTimelineParams)
		//fmt.Printf("res: %+v\n", tweets)
		if err != nil {
			log.Fatalf("err: %+v\n", err)
		}
		for _, tw := range tweets {
			fmt.Printf("tweet: %s\n  tweetID: %s\n", tw.Text, tw.IDStr)
			lastTweetID = tw.ID
			date, err := time.Parse(time.RubyDate, tw.CreatedAt)
			if err != nil {
				log.Fatalf("err: %+v\n", err)
			}
			//fmt.Printf("show dates: %s", date.Format(time.RubyDate))
			if beforeFromDate.After(date){
				fmt.Printf("  delete tweet: %s .\n", tw.Text)
				client.Statuses.Destroy(tw.ID, &twitter.StatusDestroyParams{
					ID: tw.ID,
				})
				continue
			}
			fmt.Printf("  not deleted. \n")
		}
	}
}

// delete favorites by given number
func deleteFavorites(client *twitter.Client, twitterID string, deleteFavs int) error {
	if twitterID == "" {
		return errors.New("twitterID cannot be empty")
	}

	params := twitter.FavoriteListParams{
		ScreenName: twitterID,
		Count:      200,
	}
	count := 0
	for ; count < deleteFavs; count += 200 {
		favs, _, err := client.Favorites.List(&params)
		if err != nil {
			fmt.Printf("err: %s \n", err)
		}
		fmt.Printf("res: %+v \n", favs)

		for _, t := range favs {
			tw, _, err := client.Favorites.Destroy(&twitter.FavoriteDestroyParams{
				ID: t.ID,
			})
			if err != nil {
				fmt.Printf("err; %+v \n", err)
				return err
			}
			fmt.Printf("tweet: %s, \n  deleted\n", tw.Text)
		}
		fmt.Printf("favs deleted \n")
	}
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file.")
	}
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	deleteFavorites(client, "twitterID", 1000)
}

// other examples
//params := &twitter.UserShowParams{
//	ScreenName: "yoiwki",
//}
//user, _, err := client.Users.Show(params)

//tweet, _, err := client.Statuses.Update("tweet from go-twitter", nil)
//if err != nil {
//	fmt.Printf("err: %s", err)
//}
//fmt.Printf("res: %s", tweet.Text)