package news

import "encoding/xml"

type NewListInformation struct {
	XMLName             xml.Name `xml:"NewListInformation"`
	Text                string   `xml:",chardata"`
	ClubName            string   `xml:"ClubName"`
	ClubWebsiteURL      string   `xml:"ClubWebsiteURL"`
	NewsletterNewsItems struct {
		Text               string `xml:",chardata"`
		NewsletterNewsItem []struct {
			Text              string `xml:",chardata"`
			ArticleURL        string `xml:"ArticleURL"`
			NewsArticleID     string `xml:"NewsArticleID"`
			PublishDate       string `xml:"PublishDate"`
			Taxonomies        string `xml:"Taxonomies"`
			TeaserText        string `xml:"TeaserText"`
			ThumbnailImageURL string `xml:"ThumbnailImageURL"`
			Title             string `xml:"Title"`
			OptaMatchId       string `xml:"OptaMatchId"`
			LastUpdateDate    string `xml:"LastUpdateDate"`
			IsPublished       string `xml:"IsPublished"`
		} `xml:"NewsletterNewsItem"`
	} `xml:"NewsletterNewsItems"`
}

type Data struct {
	Id          string      `json:"id"  bson:"id"`
	TeamId      string      `json:"teamId"  bson:"teamID"`
	OptaMatchId interface{} `json:"optaMatchId"  bson:"optaMatchID"`
	Title       string      `json:"title"  bson:"title"`
	Type        []string    `json:"type"  bson:"type"`
	Teaser      interface{} `json:"teaser"  bson:"teaser"`
	Content     string      `json:"content" bson:"content"`
	Url         string      `json:"url"  bson:"url"`
	ImageUrl    string      `json:"imageUrl"  bson:"imageUrl"`
	GalleryUrls interface{} `json:"galleryUrls" bson:"galleryUrls"`
	VideoUrl    interface{} `json:"videoUrl"  bson:"videoUrl"`
	Published   string      `json:"published"  bson:"published"`
}

type NewsArticle struct {
	Data     Data     `json:"data" bson:"data"`
	Metadata Metadata `json:"metadata"`
	Status   string   `json:"status" bson:"status"`
}

type Metadata struct {
	CreatedAt  string `json:"createdAt" bson:"createdAt"`
	Sort       string `json:"sort" bson:"sort"`
	TotalItems int    `json:"totalItems" bson:"totalItems"`
}
