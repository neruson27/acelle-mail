package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	dbRef struct {
		Ref interface{}        `bson:"$ref"`
		ID  primitive.ObjectID `bson:"$id"`
	}
	geoIP struct {
		Latitude     string `bson:"latitude"`
		Longitud     string `bson:"longitud"`
		Country      string `bson:"country"`
		State        string `bson:"state"`
		City         string `bson:"city"`
		ZipCode      string `bson:"zipCode"`
		Serv         string `bson:"serv"`
		IpId         string `bson:"ipid"`
		Organization string `bson:"organization"`
	}

	Contact struct {
		ID            primitive.ObjectID `bson:"_id" json:"-"`
		Email         string             `bson:"email,omitempty" json:"EMAIL,omitempty"`
		StatusLabel   string             `bson:"statusLabel,omitempty" json:"TAG,omitempty"`
		FirstName     string             `bson:"firstName,omitempty" json:"FIRST_NAME,omitempty"`
		LastName      string             `bson:"lastName,omitempty" json:"LAST_NAME,omitempty"`
		UTMS          map[string]string  `bson:"utms"`
		UtmsSource    string             `json:"UTMS_SOURCE,omitempty"`
		UtmsMedium    string             `json:"UTMS_MEDIUM,omitempty"`
		UtmsCampaign  string             `json:"UTM_CAMPAIGN,omitempty"`
		AssignedToRef *dbRef             `bson:"assignedTo,omitempty" `
		SubStatus     string             `bson:"subStatus,omitempty"`
		GeoIP         *geoIP             `bson:"geoIP,omitempty"`
	}
)
