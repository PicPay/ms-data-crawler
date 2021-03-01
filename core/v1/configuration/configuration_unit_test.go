package configuration

import (
	"context"
	v1 "github.com/PicPay/ms-data-crawler/core/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"

	"github.com/PicPay/ms-data-crawler/core/v1/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

//Struct that defines the objects that will be used in the test
type ServiceSuite struct {
	suite.Suite
	*require.Assertions

	ctrl           *gomock.Controller
	mockGateway    *MockGateway
	mockHttpClient *mocks.MockHttpClient

	service *Service
}

var (
	headers    = map[string]string{"teste": "123"}
	dateMock   = time.Date(2020, 10, 03, 14, 18, 32, 1000, time.UTC)
	identifier = "recent-itens"
	StatusOK   = 200
)

//Starting the test suite
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

//This method setup the initial configurations for the tests, this is run before all tests
func (s *ServiceSuite) SetupTest() {
	s.Assertions = require.New(s.T())

	s.ctrl = gomock.NewController(s.T())
	s.mockGateway = NewMockGateway(s.ctrl)
	s.mockHttpClient = mocks.NewMockHttpClient(s.ctrl)

	s.service = NewService(s.mockGateway, s.mockHttpClient)
}

//This method is ran after all tests have runned, it cleans the suite
func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceSuite) TestFetchConfiguration() {
	ctx := context.Background()
	s.mockGateway.EXPECT().Find(ctx, gomock.Any()).Return(getMockedConfiguration(), nil).Times(1)
	getSource := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpSourceResponse(), nil).Times(1)
	getFirstCrawler := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpCrawlerFirstResponse(), nil).Times(1)
	getSecondCrawler := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpCrawlerSecondResponse(), nil).Times(1)

	gomock.InOrder(
		getSource,
		getFirstCrawler,
		getSecondCrawler,
	)

	actualResponse, err := s.service.Fetch(ctx, identifier, headers)

	s.NoError(err)
	expectedResponse := getConfigurationResponse()

	s.Equal(expectedResponse, actualResponse)
}

func (s *ServiceSuite) TestFetchWrongConfiguration() {
	ctx := context.Background()
	s.mockGateway.EXPECT().Find(ctx, gomock.Any()).Return(getMockedConfiguration(), nil).Times(1)
	getSource := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpSourceResponse(), nil).Times(1)
	getFirstCrawler := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpCrawlerFirstResponse(), nil).Times(1)
	getSecondCrawler := s.mockHttpClient.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(getHttpCrawlerSecondResponse(), nil).Times(1)

	gomock.InOrder(
		getSource,
		getFirstCrawler,
		getSecondCrawler,
	)

	actualResponse, err := s.service.Fetch(ctx, identifier, headers)

	s.NoError(err)
	expectedResponse := getConfigurationResponse()

	s.Equal(expectedResponse, actualResponse)
}

func getMockedConfiguration() *Configuration {

	id, _ := primitive.ObjectIDFromHex("5f4fc94835dacb888fff5878")

	var configurationMocked = Configuration{
		ID:         id,
		Identifier: identifier,
		Source:     getSource(),
		Crawler:    getCrawler(),
		CreatedAt:  dateMock,
		UpdatedAt:  dateMock,
	}
	return &configurationMocked
}

func getHttpSourceResponse() v1.HttpResponse {
	response := "{\"_id\":\"7620665\",\"consumer_id\":\"7620665\",\"items\":[{\"id\":\"5f6b4a1b491a7d014310a754\",\"recency_score\":1,\"type\":\"legacy\"},{\"id\":\"5fb55bfec0e26e0e1501852a\",\"recency_score\":0.16666666666666666,\"type\":\"legacy\"}]}"
	return v1.HttpResponse{
		Status: StatusOK,
		Body:   []byte(response),
	}
}

func getHttpCrawlerFirstResponse() v1.HttpResponse {
	response := "{\"_id\":\"5fb55bfec0e26e0e1501852a\",\"available_in_ddds\":\"\",\"banner_img_url\":\"https:\\/\\/s3.amazonaws.com\\/cdn.picpay.com\\/picpay\\/sellers\\/ifood-banner.png\",\"category\":\"services\",\"created_at\":\"2020-11-18 15:38:06\",\"description\":\"Pede um xbox!\",\"description_large\":\"Compre créditos para usar no Xbox\",\"description_large_markdown\":\"Bem-vindo ao Xbox! \",\"disclaimer_markdown\":\"**Importante**: Você receberá o código PIN no valor da recarga escolhida após o pagamento. O crédito não é reembolsável pelo PicPay. <br><br>**Instruções de Resgate** <br>1. Abra seu iFood e clique em <i>Perfil <\\/i><br>2. Acesse sua <i>Carteira<\\/i> e pressione <i>Resgatar iFood Card<\\/i> <br>3. Digite ou copie e cole o código do seu iFood Card <br>4. O saldo do iFood Card estará na sua conta para ser utilizado.<br><br>**Validade dos créditos**<br>90 dias após o resgate no Ifood. \",\"enabled\":true,\"identifier\":\"xboxsx\",\"image_url\":\"https:\\/\\/s3.amazonaws.com\\/cdn.picpay.com\\/picpay\\/sellers\\/ifood-logo.png\",\"info_url\":\"https:\\/\\/cdn.picpay.com\\/picpay\\/sellers\\/ifood-terms.html\",\"min_version_android_integer\":101010,\"min_version_android_string\":\"10.10.10\",\"min_version_ios_integer\":101010,\"min_version_ios_string\":\"10.10.10\",\"name\":\"iFood\",\"offline\":false,\"offline_message\":\"Xbox offline\",\"operator\":\"xboxsx\",\"order\":0,\"sales_ranking\":0,\"selection_type\":\"list\",\"seller_id\":5,\"service\":\"digitalcodes\",\"show_in_search\":true,\"show_in_store\":true,\"spend_limit\":3000,\"student_account\":false,\"supplier\":\"5f6b4a19491a7d014310a665\",\"updated_at\":\"2020-11-18 15:38:06\",\"tags\":null,\"universal_link_android\":\"\",\"universal_link_ios\":\"\",\"universal_link_title\":\"\",\"cupom_fields\":null,\"screen_identifier\":\"5fb534fccb21457e84c5d694\"}"
	return v1.HttpResponse{
		Status: StatusOK,
		Body:   []byte(response),
	}
}

func getHttpCrawlerSecondResponse() v1.HttpResponse {
	response := "{\\\"_id\\\":\\\"5f6b4a1b491a7d014310a754\\\",\\\"name\\\":\\\"Cliente Premium Herbalife Nutrition\\\",\\\"description\\\":\\\"Adquira sua conta e tenha acesso a diversos benefícios\\\",\\\"service\\\":\\\"digitalcodes\\\",\\\"enabled\\\":true,\\\"seller_id\\\":14427,\\\"image_url\\\":\\\"https:\\\\/\\\\/www.picpay.com\\\\/static\\\\/images\\\\/bullets\\\\/store.png\\\",\\\"banner_img_url\\\":\\\"https:\\\\/\\\\/pngimage.net\\\\/wp-content\\\\/uploads\\\\/2018\\\\/06\\\\/picpay-png-3.png\\\",\\\"info_url\\\":\\\"https:\\\\/\\\\/cdn.picpay.com\\\\/apps\\\\/picpay\\\\/digital-goods\\\\/herbalife-terms.html\\\",\\\"description_large\\\":\\\"\\\",\\\"description_large_markdown\\\":\\\"Agora você pode adquirir produtos Herbalife Nutrition de maneira prática e rápida, e ainda garantir descontos especiais e promoções exclusivas!\\\",\\\"disclaimer_markdown\\\":\\\"**Importante:** Você receberá o ID e código de acesso ao Portal Herbalife Nutrition após o pagamento.\\\",\\\"min_version_android_string\\\":\\\"9.5.12\\\",\\\"min_version_android_integer\\\":905012,\\\"min_version_ios_string\\\":\\\"100.12.36\\\",\\\"min_version_ios_integer\\\":10012036,\\\"selection_type\\\":\\\"list\\\",\\\"order\\\":21,\\\"identifier\\\":\\\"herbalife\\\",\\\"category\\\":\\\"services\\\",\\\"sales_ranking\\\":0,\\\"operator\\\":\\\"\\\",\\\"spend_limit\\\":0,\\\"available_in_ddds\\\":\\\"\\\",\\\"offline\\\":false,\\\"offline_message\\\":\\\"Serviço indisponível no momento, por favor tente novamente mais tarde\\\",\\\"supplier\\\":\\\"herbalife\\\",\\\"show_in_store\\\":false,\\\"show_in_search\\\":false,\\\"student_account\\\":false,\\\"updated_at\\\":\\\"2020-09-23 10:14:03\\\",\\\"created_at\\\":\\\"2020-09-23 10:14:03\\\"}"
	return v1.HttpResponse{
		Status: StatusOK,
		Body:   []byte(response),
	}
}

func getSource() ServiceRequest {
	return ServiceRequest{
		Name:           "midgard",
		Url:            "http://midgard.ms.prod/v1/collections/store_recently_viewed/<Consumer_Id>",
		Method:         "GET",
		ContentType:    "application/json",
		Validation:     []KeyValue{},
		ResponseFormat: "{{if .items}}\n{{$total := len .items}}\n  [\n  {{ range $key, $value := .items }}\n  {\n    \"id\": \"{{$value.id}}\",\n    \"type\": \"{{$value.type}}\"\n  }\n  {{if HasMoreItems $key $total}}\n   ,\n  {{end}}\n  {{end}}\n  ]\n{{end}}",
	}
}

func getCrawler() []ServiceRequest {
	return []ServiceRequest{
		{
			Name:        "digital-goods",
			Url:         "http://store.sandbox.limbo.work:8180/v1/digitalgoods/services/:dgId",
			Headers:     []KeyValue{},
			Method:      "GET",
			ContentType: "application/json",
			Mapping: []KeyValue{
				{Index: ":dgId", Value: "id"},
			},
			Validation: []KeyValue{
				{Index: "type", Value: "legacy"},
			},
		},
	}
}

func getConfigurationResponse() *AssembledScreen {
	return &AssembledScreen{
		Identifier: identifier,
		Data: []interface{}{
			map[string]interface{}{"_id": "5fb55bfec0e26e0e1501852a", "available_in_ddds": "", "banner_img_url": "https://s3.amazonaws.com/cdn.picpay.com/picpay/sellers/ifood-banner.png", "category": "services", "created_at": "2020-11-18 15:38:06", "cupom_fields": nil, "description": "Pede um xbox!", "description_large": "Compre créditos para usar no Xbox", "description_large_markdown": "Bem-vindo ao Xbox! ", "disclaimer_markdown": "**Importante**: Você receberá o código PIN no valor da recarga escolhida após o pagamento. O crédito não é reembolsável pelo PicPay. \u003cbr\u003e\u003cbr\u003e**Instruções de Resgate** \u003cbr\u003e1. Abra seu iFood e clique em \u003ci\u003ePerfil \u003c/i\u003e\u003cbr\u003e2. Acesse sua \u003ci\u003eCarteira\u003c/i\u003e e pressione \u003ci\u003eResgatar iFood Card\u003c/i\u003e \u003cbr\u003e3. Digite ou copie e cole o código do seu iFood Card \u003cbr\u003e4. O saldo do iFood Card estará na sua conta para ser utilizado.\u003cbr\u003e\u003cbr\u003e**Validade dos créditos**\u003cbr\u003e90 dias após o resgate no Ifood. ", "enabled": true, "identifier": "xboxsx", "image_url": "https://s3.amazonaws.com/cdn.picpay.com/picpay/sellers/ifood-logo.png", "info_url": "https://cdn.picpay.com/picpay/sellers/ifood-terms.html", "min_version_android_integer": float64(101010), "min_version_android_string": "10.10.10", "min_version_ios_integer": float64(101010), "min_version_ios_string": "10.10.10", "name": "iFood", "offline": false, "offline_message": "Xbox offline", "operator": "xboxsx", "order": float64(0), "sales_ranking": float64(0), "screen_identifier": "5fb534fccb21457e84c5d694", "selection_type": "list", "seller_id": float64(5), "service": "digitalcodes", "show_in_search": true, "show_in_store": true, "spend_limit": float64(3000), "student_account": false, "supplier": "5f6b4a19491a7d014310a665", "tags": nil, "universal_link_android": "", "universal_link_ios": "", "universal_link_title": "", "updated_at": "2020-11-18 15:38:06"},
		},
		CreatedAt: dateMock,
		UpdatedAt: dateMock,
	}
}
