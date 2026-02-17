// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// socialLoginFields maps TF map keys to SDK User struct getters/setters.
var socialLoginFields = []struct {
	Key string
	Get func(*casdoorsdk.User) string
	Set func(*casdoorsdk.User, string)
}{
	{"github", func(u *casdoorsdk.User) string { return u.GitHub }, func(u *casdoorsdk.User, v string) { u.GitHub = v }},
	{"google", func(u *casdoorsdk.User) string { return u.Google }, func(u *casdoorsdk.User, v string) { u.Google = v }},
	{"qq", func(u *casdoorsdk.User) string { return u.QQ }, func(u *casdoorsdk.User, v string) { u.QQ = v }},
	{"wechat", func(u *casdoorsdk.User) string { return u.WeChat }, func(u *casdoorsdk.User, v string) { u.WeChat = v }},
	{"facebook", func(u *casdoorsdk.User) string { return u.Facebook }, func(u *casdoorsdk.User, v string) { u.Facebook = v }},
	{"dingtalk", func(u *casdoorsdk.User) string { return u.DingTalk }, func(u *casdoorsdk.User, v string) { u.DingTalk = v }},
	{"weibo", func(u *casdoorsdk.User) string { return u.Weibo }, func(u *casdoorsdk.User, v string) { u.Weibo = v }},
	{"gitee", func(u *casdoorsdk.User) string { return u.Gitee }, func(u *casdoorsdk.User, v string) { u.Gitee = v }},
	{"linkedin", func(u *casdoorsdk.User) string { return u.LinkedIn }, func(u *casdoorsdk.User, v string) { u.LinkedIn = v }},
	{"wecom", func(u *casdoorsdk.User) string { return u.Wecom }, func(u *casdoorsdk.User, v string) { u.Wecom = v }},
	{"lark", func(u *casdoorsdk.User) string { return u.Lark }, func(u *casdoorsdk.User, v string) { u.Lark = v }},
	{"gitlab", func(u *casdoorsdk.User) string { return u.Gitlab }, func(u *casdoorsdk.User, v string) { u.Gitlab = v }},
	{"adfs", func(u *casdoorsdk.User) string { return u.Adfs }, func(u *casdoorsdk.User, v string) { u.Adfs = v }},
	{"baidu", func(u *casdoorsdk.User) string { return u.Baidu }, func(u *casdoorsdk.User, v string) { u.Baidu = v }},
	{"alipay", func(u *casdoorsdk.User) string { return u.Alipay }, func(u *casdoorsdk.User, v string) { u.Alipay = v }},
	{"casdoor", func(u *casdoorsdk.User) string { return u.Casdoor }, func(u *casdoorsdk.User, v string) { u.Casdoor = v }},
	{"infoflow", func(u *casdoorsdk.User) string { return u.Infoflow }, func(u *casdoorsdk.User, v string) { u.Infoflow = v }},
	{"apple", func(u *casdoorsdk.User) string { return u.Apple }, func(u *casdoorsdk.User, v string) { u.Apple = v }},
	{"azuread", func(u *casdoorsdk.User) string { return u.AzureAD }, func(u *casdoorsdk.User, v string) { u.AzureAD = v }},
	{"azureadb2c", func(u *casdoorsdk.User) string { return u.AzureADB2c }, func(u *casdoorsdk.User, v string) { u.AzureADB2c = v }},
	{"slack", func(u *casdoorsdk.User) string { return u.Slack }, func(u *casdoorsdk.User, v string) { u.Slack = v }},
	{"steam", func(u *casdoorsdk.User) string { return u.Steam }, func(u *casdoorsdk.User, v string) { u.Steam = v }},
	{"bilibili", func(u *casdoorsdk.User) string { return u.Bilibili }, func(u *casdoorsdk.User, v string) { u.Bilibili = v }},
	{"okta", func(u *casdoorsdk.User) string { return u.Okta }, func(u *casdoorsdk.User, v string) { u.Okta = v }},
	{"douyin", func(u *casdoorsdk.User) string { return u.Douyin }, func(u *casdoorsdk.User, v string) { u.Douyin = v }},
	{"kwai", func(u *casdoorsdk.User) string { return u.Kwai }, func(u *casdoorsdk.User, v string) { u.Kwai = v }},
	{"line", func(u *casdoorsdk.User) string { return u.Line }, func(u *casdoorsdk.User, v string) { u.Line = v }},
	{"amazon", func(u *casdoorsdk.User) string { return u.Amazon }, func(u *casdoorsdk.User, v string) { u.Amazon = v }},
	{"auth0", func(u *casdoorsdk.User) string { return u.Auth0 }, func(u *casdoorsdk.User, v string) { u.Auth0 = v }},
	{"battlenet", func(u *casdoorsdk.User) string { return u.BattleNet }, func(u *casdoorsdk.User, v string) { u.BattleNet = v }},
	{"bitbucket", func(u *casdoorsdk.User) string { return u.Bitbucket }, func(u *casdoorsdk.User, v string) { u.Bitbucket = v }},
	{"box", func(u *casdoorsdk.User) string { return u.Box }, func(u *casdoorsdk.User, v string) { u.Box = v }},
	{"cloudfoundry", func(u *casdoorsdk.User) string { return u.CloudFoundry }, func(u *casdoorsdk.User, v string) { u.CloudFoundry = v }},
	{"dailymotion", func(u *casdoorsdk.User) string { return u.Dailymotion }, func(u *casdoorsdk.User, v string) { u.Dailymotion = v }},
	{"deezer", func(u *casdoorsdk.User) string { return u.Deezer }, func(u *casdoorsdk.User, v string) { u.Deezer = v }},
	{"digitalocean", func(u *casdoorsdk.User) string { return u.DigitalOcean }, func(u *casdoorsdk.User, v string) { u.DigitalOcean = v }},
	{"discord", func(u *casdoorsdk.User) string { return u.Discord }, func(u *casdoorsdk.User, v string) { u.Discord = v }},
	{"dropbox", func(u *casdoorsdk.User) string { return u.Dropbox }, func(u *casdoorsdk.User, v string) { u.Dropbox = v }},
	{"eveonline", func(u *casdoorsdk.User) string { return u.EveOnline }, func(u *casdoorsdk.User, v string) { u.EveOnline = v }},
	{"fitbit", func(u *casdoorsdk.User) string { return u.Fitbit }, func(u *casdoorsdk.User, v string) { u.Fitbit = v }},
	{"gitea", func(u *casdoorsdk.User) string { return u.Gitea }, func(u *casdoorsdk.User, v string) { u.Gitea = v }},
	{"heroku", func(u *casdoorsdk.User) string { return u.Heroku }, func(u *casdoorsdk.User, v string) { u.Heroku = v }},
	{"influxcloud", func(u *casdoorsdk.User) string { return u.InfluxCloud }, func(u *casdoorsdk.User, v string) { u.InfluxCloud = v }},
	{"instagram", func(u *casdoorsdk.User) string { return u.Instagram }, func(u *casdoorsdk.User, v string) { u.Instagram = v }},
	{"intercom", func(u *casdoorsdk.User) string { return u.Intercom }, func(u *casdoorsdk.User, v string) { u.Intercom = v }},
	{"kakao", func(u *casdoorsdk.User) string { return u.Kakao }, func(u *casdoorsdk.User, v string) { u.Kakao = v }},
	{"lastfm", func(u *casdoorsdk.User) string { return u.Lastfm }, func(u *casdoorsdk.User, v string) { u.Lastfm = v }},
	{"mailru", func(u *casdoorsdk.User) string { return u.Mailru }, func(u *casdoorsdk.User, v string) { u.Mailru = v }},
	{"meetup", func(u *casdoorsdk.User) string { return u.Meetup }, func(u *casdoorsdk.User, v string) { u.Meetup = v }},
	{"microsoftonline", func(u *casdoorsdk.User) string { return u.MicrosoftOnline }, func(u *casdoorsdk.User, v string) { u.MicrosoftOnline = v }},
	{"naver", func(u *casdoorsdk.User) string { return u.Naver }, func(u *casdoorsdk.User, v string) { u.Naver = v }},
	{"nextcloud", func(u *casdoorsdk.User) string { return u.Nextcloud }, func(u *casdoorsdk.User, v string) { u.Nextcloud = v }},
	{"onedrive", func(u *casdoorsdk.User) string { return u.OneDrive }, func(u *casdoorsdk.User, v string) { u.OneDrive = v }},
	{"oura", func(u *casdoorsdk.User) string { return u.Oura }, func(u *casdoorsdk.User, v string) { u.Oura = v }},
	{"patreon", func(u *casdoorsdk.User) string { return u.Patreon }, func(u *casdoorsdk.User, v string) { u.Patreon = v }},
	{"paypal", func(u *casdoorsdk.User) string { return u.Paypal }, func(u *casdoorsdk.User, v string) { u.Paypal = v }},
	{"salesforce", func(u *casdoorsdk.User) string { return u.SalesForce }, func(u *casdoorsdk.User, v string) { u.SalesForce = v }},
	{"shopify", func(u *casdoorsdk.User) string { return u.Shopify }, func(u *casdoorsdk.User, v string) { u.Shopify = v }},
	{"soundcloud", func(u *casdoorsdk.User) string { return u.Soundcloud }, func(u *casdoorsdk.User, v string) { u.Soundcloud = v }},
	{"spotify", func(u *casdoorsdk.User) string { return u.Spotify }, func(u *casdoorsdk.User, v string) { u.Spotify = v }},
	{"strava", func(u *casdoorsdk.User) string { return u.Strava }, func(u *casdoorsdk.User, v string) { u.Strava = v }},
	{"stripe", func(u *casdoorsdk.User) string { return u.Stripe }, func(u *casdoorsdk.User, v string) { u.Stripe = v }},
	{"tiktok", func(u *casdoorsdk.User) string { return u.TikTok }, func(u *casdoorsdk.User, v string) { u.TikTok = v }},
	{"tumblr", func(u *casdoorsdk.User) string { return u.Tumblr }, func(u *casdoorsdk.User, v string) { u.Tumblr = v }},
	{"twitch", func(u *casdoorsdk.User) string { return u.Twitch }, func(u *casdoorsdk.User, v string) { u.Twitch = v }},
	{"twitter", func(u *casdoorsdk.User) string { return u.Twitter }, func(u *casdoorsdk.User, v string) { u.Twitter = v }},
	{"typetalk", func(u *casdoorsdk.User) string { return u.Typetalk }, func(u *casdoorsdk.User, v string) { u.Typetalk = v }},
	{"uber", func(u *casdoorsdk.User) string { return u.Uber }, func(u *casdoorsdk.User, v string) { u.Uber = v }},
	{"vk", func(u *casdoorsdk.User) string { return u.VK }, func(u *casdoorsdk.User, v string) { u.VK = v }},
	{"wepay", func(u *casdoorsdk.User) string { return u.Wepay }, func(u *casdoorsdk.User, v string) { u.Wepay = v }},
	{"xero", func(u *casdoorsdk.User) string { return u.Xero }, func(u *casdoorsdk.User, v string) { u.Xero = v }},
	{"yahoo", func(u *casdoorsdk.User) string { return u.Yahoo }, func(u *casdoorsdk.User, v string) { u.Yahoo = v }},
	{"yammer", func(u *casdoorsdk.User) string { return u.Yammer }, func(u *casdoorsdk.User, v string) { u.Yammer = v }},
	{"yandex", func(u *casdoorsdk.User) string { return u.Yandex }, func(u *casdoorsdk.User, v string) { u.Yandex = v }},
	{"zoom", func(u *casdoorsdk.User) string { return u.Zoom }, func(u *casdoorsdk.User, v string) { u.Zoom = v }},
	{"metamask", func(u *casdoorsdk.User) string { return u.MetaMask }, func(u *casdoorsdk.User, v string) { u.MetaMask = v }},
	{"web3onboard", func(u *casdoorsdk.User) string { return u.Web3Onboard }, func(u *casdoorsdk.User, v string) { u.Web3Onboard = v }},
	{"custom", func(u *casdoorsdk.User) string { return u.Custom }, func(u *casdoorsdk.User, v string) { u.Custom = v }},
	{"custom2", func(u *casdoorsdk.User) string { return u.Custom2 }, func(u *casdoorsdk.User, v string) { u.Custom2 = v }},
	{"custom3", func(u *casdoorsdk.User) string { return u.Custom3 }, func(u *casdoorsdk.User, v string) { u.Custom3 = v }},
	{"custom4", func(u *casdoorsdk.User) string { return u.Custom4 }, func(u *casdoorsdk.User, v string) { u.Custom4 = v }},
	{"custom5", func(u *casdoorsdk.User) string { return u.Custom5 }, func(u *casdoorsdk.User, v string) { u.Custom5 = v }},
	{"custom6", func(u *casdoorsdk.User) string { return u.Custom6 }, func(u *casdoorsdk.User, v string) { u.Custom6 = v }},
	{"custom7", func(u *casdoorsdk.User) string { return u.Custom7 }, func(u *casdoorsdk.User, v string) { u.Custom7 = v }},
	{"custom8", func(u *casdoorsdk.User) string { return u.Custom8 }, func(u *casdoorsdk.User, v string) { u.Custom8 = v }},
	{"custom9", func(u *casdoorsdk.User) string { return u.Custom9 }, func(u *casdoorsdk.User, v string) { u.Custom9 = v }},
	{"custom10", func(u *casdoorsdk.User) string { return u.Custom10 }, func(u *casdoorsdk.User, v string) { u.Custom10 = v }},
}

// AddressModel represents a structured address entry.
type AddressModel struct {
	Tag     types.String `tfsdk:"tag"`
	Line1   types.String `tfsdk:"line1"`
	Line2   types.String `tfsdk:"line2"`
	City    types.String `tfsdk:"city"`
	State   types.String `tfsdk:"state"`
	ZipCode types.String `tfsdk:"zip_code"`
	Region  types.String `tfsdk:"region"`
}

// AddressAttrTypes returns the attribute types for AddressModel.
func AddressAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag":      types.StringType,
		"line1":    types.StringType,
		"line2":    types.StringType,
		"city":     types.StringType,
		"state":    types.StringType,
		"zip_code": types.StringType,
		"region":   types.StringType,
	}
}

// ManagedAccountModel represents a managed account entry.
type ManagedAccountModel struct {
	Application types.String `tfsdk:"application"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	SigninUrl   types.String `tfsdk:"signin_url"`
}

// ManagedAccountAttrTypes returns the attribute types for ManagedAccountModel.
func ManagedAccountAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"application": types.StringType,
		"username":    types.StringType,
		"password":    types.StringType,
		"signin_url":  types.StringType,
	}
}

// MfaAccountModel represents an MFA account entry.
type MfaAccountModel struct {
	AccountName types.String `tfsdk:"account_name"`
	Issuer      types.String `tfsdk:"issuer"`
	SecretKey   types.String `tfsdk:"secret_key"`
	Origin      types.String `tfsdk:"origin"`
}

// MfaAccountAttrTypes returns the attribute types for MfaAccountModel.
func MfaAccountAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_name": types.StringType,
		"issuer":       types.StringType,
		"secret_key":   types.StringType,
		"origin":       types.StringType,
	}
}

// FaceIdModel represents a face ID entry.
type FaceIdModel struct {
	Name       types.String `tfsdk:"name"`
	FaceIdData types.List   `tfsdk:"face_id_data"`
	ImageUrl   types.String `tfsdk:"image_url"`
}

// FaceIdAttrTypes returns the attribute types for FaceIdModel.
func FaceIdAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"face_id_data": types.ListType{ElemType: types.Float64Type},
		"image_url":    types.StringType,
	}
}

// ProductInfoModel represents a cart product info entry.
type ProductInfoModel struct {
	Owner       types.String  `tfsdk:"owner"`
	Name        types.String  `tfsdk:"name"`
	DisplayName types.String  `tfsdk:"display_name"`
	Image       types.String  `tfsdk:"image"`
	Detail      types.String  `tfsdk:"detail"`
	Price       types.Float64 `tfsdk:"price"`
	Currency    types.String  `tfsdk:"currency"`
	IsRecharge  types.Bool    `tfsdk:"is_recharge"`
	Quantity    types.Int64   `tfsdk:"quantity"`
	PricingName types.String  `tfsdk:"pricing_name"`
	PlanName    types.String  `tfsdk:"plan_name"`
}

// ProductInfoAttrTypes returns the attribute types for ProductInfoModel.
func ProductInfoAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"owner":        types.StringType,
		"name":         types.StringType,
		"display_name": types.StringType,
		"image":        types.StringType,
		"detail":       types.StringType,
		"price":        types.Float64Type,
		"currency":     types.StringType,
		"is_recharge":  types.BoolType,
		"quantity":     types.Int64Type,
		"pricing_name": types.StringType,
		"plan_name":    types.StringType,
	}
}

type UserResource struct {
	client *casdoorsdk.Client
}

type UserResourceModel struct {
	Owner                  types.String  `tfsdk:"owner"`
	Name                   types.String  `tfsdk:"name"`
	ID                     types.String  `tfsdk:"id"`
	Type                   types.String  `tfsdk:"type"`
	Password               types.String  `tfsdk:"password"`
	PasswordType           types.String  `tfsdk:"password_type"`
	DisplayName            types.String  `tfsdk:"display_name"`
	FirstName              types.String  `tfsdk:"first_name"`
	LastName               types.String  `tfsdk:"last_name"`
	Avatar                 types.String  `tfsdk:"avatar"`
	Email                  types.String  `tfsdk:"email"`
	EmailVerified          types.Bool    `tfsdk:"email_verified"`
	Phone                  types.String  `tfsdk:"phone"`
	CountryCode            types.String  `tfsdk:"country_code"`
	Region                 types.String  `tfsdk:"region"`
	Location               types.String  `tfsdk:"location"`
	Affiliation            types.String  `tfsdk:"affiliation"`
	Title                  types.String  `tfsdk:"title"`
	Homepage               types.String  `tfsdk:"homepage"`
	Bio                    types.String  `tfsdk:"bio"`
	Tag                    types.String  `tfsdk:"tag"`
	Language               types.String  `tfsdk:"language"`
	Gender                 types.String  `tfsdk:"gender"`
	Birthday               types.String  `tfsdk:"birthday"`
	Education              types.String  `tfsdk:"education"`
	Score                  types.Int64   `tfsdk:"score"`
	Karma                  types.Int64   `tfsdk:"karma"`
	Ranking                types.Int64   `tfsdk:"ranking"`
	IsAdmin                types.Bool    `tfsdk:"is_admin"`
	IsForbidden            types.Bool    `tfsdk:"is_forbidden"`
	IsDeleted              types.Bool    `tfsdk:"is_deleted"`
	SignupApplication      types.String  `tfsdk:"signup_application"`
	CreatedTime            types.String  `tfsdk:"created_time"`
	UpdatedTime            types.String  `tfsdk:"updated_time"`
	DeletedTime            types.String  `tfsdk:"deleted_time"`
	ExternalId             types.String  `tfsdk:"external_id"`
	PasswordSalt           types.String  `tfsdk:"password_salt"`
	AvatarType             types.String  `tfsdk:"avatar_type"`
	PermanentAvatar        types.String  `tfsdk:"permanent_avatar"`
	Address                types.List    `tfsdk:"address"`
	Addresses              types.List    `tfsdk:"addresses"`
	IdCardType             types.String  `tfsdk:"id_card_type"`
	IdCard                 types.String  `tfsdk:"id_card"`
	RealName               types.String  `tfsdk:"real_name"`
	IsVerified             types.Bool    `tfsdk:"is_verified"`
	IsDefaultAvatar        types.Bool    `tfsdk:"is_default_avatar"`
	IsOnline               types.Bool    `tfsdk:"is_online"`
	Hash                   types.String  `tfsdk:"hash"`
	PreHash                types.String  `tfsdk:"pre_hash"`
	Balance                types.Float64 `tfsdk:"balance"`
	BalanceCredit          types.Float64 `tfsdk:"balance_credit"`
	Currency               types.String  `tfsdk:"currency"`
	BalanceCurrency        types.String  `tfsdk:"balance_currency"`
	RegisterType           types.String  `tfsdk:"register_type"`
	RegisterSource         types.String  `tfsdk:"register_source"`
	AccessKey              types.String  `tfsdk:"access_key"`
	AccessSecret           types.String  `tfsdk:"access_secret"`
	AccessToken            types.String  `tfsdk:"access_token"`
	OriginalToken          types.String  `tfsdk:"original_token"`
	OriginalRefreshToken   types.String  `tfsdk:"original_refresh_token"`
	CreatedIp              types.String  `tfsdk:"created_ip"`
	LastSigninTime         types.String  `tfsdk:"last_signin_time"`
	LastSigninIp           types.String  `tfsdk:"last_signin_ip"`
	SocialLogins           types.Map     `tfsdk:"social_logins"`
	Invitation             types.String  `tfsdk:"invitation"`
	InvitationCode         types.String  `tfsdk:"invitation_code"`
	Ldap                   types.String  `tfsdk:"ldap"`
	Properties             types.Map     `tfsdk:"properties"`
	NeedUpdatePassword     types.Bool    `tfsdk:"need_update_password"`
	LastChangePasswordTime types.String  `tfsdk:"last_change_password_time"`
	LastSigninWrongTime    types.String  `tfsdk:"last_signin_wrong_time"`
	SigninWrongTimes       types.Int64   `tfsdk:"signin_wrong_times"`
	PreferredMfaType       types.String  `tfsdk:"preferred_mfa_type"`
	RecoveryCodes          types.List    `tfsdk:"recovery_codes"`
	TotpSecret             types.String  `tfsdk:"totp_secret"`
	MfaPhoneEnabled        types.Bool    `tfsdk:"mfa_phone_enabled"`
	MfaEmailEnabled        types.Bool    `tfsdk:"mfa_email_enabled"`
	MfaRadiusEnabled       types.Bool    `tfsdk:"mfa_radius_enabled"`
	MfaRadiusUsername      types.String  `tfsdk:"mfa_radius_username"`
	MfaRadiusProvider      types.String  `tfsdk:"mfa_radius_provider"`
	MfaPushEnabled         types.Bool    `tfsdk:"mfa_push_enabled"`
	MfaPushReceiver        types.String  `tfsdk:"mfa_push_receiver"`
	MfaPushProvider        types.String  `tfsdk:"mfa_push_provider"`
	MfaRememberDeadline    types.String  `tfsdk:"mfa_remember_deadline"`
	IpWhitelist            types.String  `tfsdk:"ip_whitelist"`
	ManagedAccounts        types.List    `tfsdk:"managed_accounts"`
	MfaAccounts            types.List    `tfsdk:"mfa_accounts"`
	MfaItems               types.List    `tfsdk:"mfa_items"`
	FaceIds                types.List    `tfsdk:"face_ids"`
	Cart                   types.List    `tfsdk:"cart"`
	Groups                 types.List    `tfsdk:"groups"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Casdoor user.",
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Description: "The organization that owns this user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The unique username.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the user in the format 'owner/name'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The user type (e.g., 'normal-user').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("normal-user"),
			},
			"password": schema.StringAttribute{
				Description: "The user's password. Note: This is write-only and will not be read back from Casdoor.",
				Optional:    true,
				Sensitive:   true,
			},
			"password_type": schema.StringAttribute{
				Description: "The password hashing type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"first_name": schema.StringAttribute{
				Description: "The user's first name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"last_name": schema.StringAttribute{
				Description: "The user's last name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"avatar": schema.StringAttribute{
				Description: "URL of the user's avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"email_verified": schema.BoolAttribute{
				Description: "Whether the user's email has been verified.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"phone": schema.StringAttribute{
				Description: "The user's phone number.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"country_code": schema.StringAttribute{
				Description: "The country code for the phone number.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"region": schema.StringAttribute{
				Description: "The user's region.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"location": schema.StringAttribute{
				Description: "The user's location.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"affiliation": schema.StringAttribute{
				Description: "The user's affiliation (e.g., company name).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"title": schema.StringAttribute{
				Description: "The user's job title.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"homepage": schema.StringAttribute{
				Description: "The user's homepage URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"bio": schema.StringAttribute{
				Description: "The user's biography.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tag": schema.StringAttribute{
				Description: "A tag for the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"language": schema.StringAttribute{
				Description: "The user's preferred language.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"gender": schema.StringAttribute{
				Description: "The user's gender.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"birthday": schema.StringAttribute{
				Description: "The user's birthday (ISO 8601 format).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"education": schema.StringAttribute{
				Description: "The user's education level.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"score": schema.Int64Attribute{
				Description: "The user's score.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"karma": schema.Int64Attribute{
				Description: "The user's karma points.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"ranking": schema.Int64Attribute{
				Description: "The user's ranking.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"is_admin": schema.BoolAttribute{
				Description: "Whether the user is an administrator.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_forbidden": schema.BoolAttribute{
				Description: "Whether the user is forbidden (disabled).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_deleted": schema.BoolAttribute{
				Description: "Whether the user is soft-deleted.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"signup_application": schema.StringAttribute{
				Description: "The application through which the user signed up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"created_time": schema.StringAttribute{
				Description: "The time when the user was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_time": schema.StringAttribute{
				Description: "The time when the user was last updated.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deleted_time": schema.StringAttribute{
				Description: "The time when the user was soft-deleted. Server-managed.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_id": schema.StringAttribute{
				Description: "External ID for the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"password_salt": schema.StringAttribute{
				Description: "The password salt.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"avatar_type": schema.StringAttribute{
				Description: "The type of avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"permanent_avatar": schema.StringAttribute{
				Description: "URL of the permanent avatar.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"address": schema.ListAttribute{
				Description: "The user's address lines.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"addresses": schema.ListNestedAttribute{
				Description: "The user's structured addresses.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag": schema.StringAttribute{
							Description: "Address tag/label.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"line1": schema.StringAttribute{
							Description: "Address line 1.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"line2": schema.StringAttribute{
							Description: "Address line 2.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"city": schema.StringAttribute{
							Description: "City.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"state": schema.StringAttribute{
							Description: "State/province.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"zip_code": schema.StringAttribute{
							Description: "ZIP/postal code.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"region": schema.StringAttribute{
							Description: "Region/country.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"id_card_type": schema.StringAttribute{
				Description: "The type of ID card.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"id_card": schema.StringAttribute{
				Description: "The ID card number.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"real_name": schema.StringAttribute{
				Description: "The user's real name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"is_verified": schema.BoolAttribute{
				Description: "Whether the user is verified.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_default_avatar": schema.BoolAttribute{
				Description: "Whether the user has the default avatar.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_online": schema.BoolAttribute{
				Description: "Whether the user is currently online.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hash": schema.StringAttribute{
				Description: "The user hash.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pre_hash": schema.StringAttribute{
				Description: "The previous user hash.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"balance": schema.Float64Attribute{
				Description: "The user's balance.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"balance_credit": schema.Float64Attribute{
				Description: "The user's balance credit.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"currency": schema.StringAttribute{
				Description: "The user's currency.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"balance_currency": schema.StringAttribute{
				Description: "The user's balance currency.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"register_type": schema.StringAttribute{
				Description: "The registration type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"register_source": schema.StringAttribute{
				Description: "The registration source.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"access_key": schema.StringAttribute{
				Description: "The user's access key.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"access_secret": schema.StringAttribute{
				Description: "The user's access secret.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"access_token": schema.StringAttribute{
				Description: "The user's access token.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"original_token": schema.StringAttribute{
				Description: "The user's original token.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"original_refresh_token": schema.StringAttribute{
				Description: "The user's original refresh token.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"created_ip": schema.StringAttribute{
				Description: "The IP address the user was created from.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_time": schema.StringAttribute{
				Description: "The last sign-in time.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_ip": schema.StringAttribute{
				Description: "The last sign-in IP address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"social_logins": schema.MapAttribute{
				Description: "Social login provider IDs. Keys are provider names (e.g., 'github', 'google').",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"invitation": schema.StringAttribute{
				Description: "The invitation used to sign up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"invitation_code": schema.StringAttribute{
				Description: "The invitation code used to sign up.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ldap": schema.StringAttribute{
				Description: "LDAP identifier.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"properties": schema.MapAttribute{
				Description: "Custom properties for the user.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"need_update_password": schema.BoolAttribute{
				Description: "Whether the user needs to update their password.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"last_change_password_time": schema.StringAttribute{
				Description: "The last time the password was changed.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_signin_wrong_time": schema.StringAttribute{
				Description: "The last time a wrong sign-in attempt was made.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"signin_wrong_times": schema.Int64Attribute{
				Description: "The number of wrong sign-in attempts.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"preferred_mfa_type": schema.StringAttribute{
				Description: "The preferred MFA type.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"recovery_codes": schema.ListAttribute{
				Description: "MFA recovery codes.",
				Optional:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"totp_secret": schema.StringAttribute{
				Description: "The TOTP secret for MFA.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_phone_enabled": schema.BoolAttribute{
				Description: "Whether phone-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mfa_email_enabled": schema.BoolAttribute{
				Description: "Whether email-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mfa_radius_enabled": schema.BoolAttribute{
				Description: "Whether RADIUS-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mfa_radius_username": schema.StringAttribute{
				Description: "The RADIUS MFA username.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_radius_provider": schema.StringAttribute{
				Description: "The RADIUS MFA provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_push_enabled": schema.BoolAttribute{
				Description: "Whether push-based MFA is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mfa_push_receiver": schema.StringAttribute{
				Description: "The push MFA receiver.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_push_provider": schema.StringAttribute{
				Description: "The push MFA provider.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"mfa_remember_deadline": schema.StringAttribute{
				Description: "The MFA remember deadline.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ip_whitelist": schema.StringAttribute{
				Description: "The IP whitelist for the user.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"managed_accounts": schema.ListNestedAttribute{
				Description: "The user's managed accounts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"application": schema.StringAttribute{
							Description: "The application name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"username": schema.StringAttribute{
							Description: "The account username.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"password": schema.StringAttribute{
							Description: "The account password.",
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Default:     stringdefault.StaticString(""),
						},
						"signin_url": schema.StringAttribute{
							Description: "The sign-in URL.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"mfa_accounts": schema.ListNestedAttribute{
				Description: "The user's MFA accounts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_name": schema.StringAttribute{
							Description: "The MFA account name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"issuer": schema.StringAttribute{
							Description: "The MFA issuer.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"secret_key": schema.StringAttribute{
							Description: "The MFA secret key.",
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Default:     stringdefault.StaticString(""),
						},
						"origin": schema.StringAttribute{
							Description: "The MFA origin.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"mfa_items": schema.ListNestedAttribute{
				Description: "The user's MFA items.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The MFA item name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"rule": schema.StringAttribute{
							Description: "The MFA item rule.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"face_ids": schema.ListNestedAttribute{
				Description: "The user's face IDs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The face ID name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"face_id_data": schema.ListAttribute{
							Description: "The face ID data points.",
							Optional:    true,
							ElementType: types.Float64Type,
						},
						"image_url": schema.StringAttribute{
							Description: "The face image URL.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"cart": schema.ListNestedAttribute{
				Description: "The user's shopping cart.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"owner": schema.StringAttribute{
							Description: "The product owner.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"name": schema.StringAttribute{
							Description: "The product name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"display_name": schema.StringAttribute{
							Description: "The product display name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"image": schema.StringAttribute{
							Description: "The product image URL.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"detail": schema.StringAttribute{
							Description: "The product detail.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"price": schema.Float64Attribute{
							Description: "The product price.",
							Optional:    true,
							Computed:    true,
							Default:     float64default.StaticFloat64(0),
						},
						"currency": schema.StringAttribute{
							Description: "The product currency.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"is_recharge": schema.BoolAttribute{
							Description: "Whether this is a recharge product.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"quantity": schema.Int64Attribute{
							Description: "The product quantity.",
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(0),
						},
						"pricing_name": schema.StringAttribute{
							Description: "The pricing name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"plan_name": schema.StringAttribute{
							Description: "The plan name.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"groups": schema.ListAttribute{
				Description: "List of groups the user belongs to.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*casdoorsdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *casdoorsdk.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// userPlanToSDK converts a UserResourceModel into a casdoorsdk.User struct.
// createdTime and updatedTime are passed explicitly because Create and Update
// handle them differently. internalID is the Casdoor-internal Id field
// (populated only during Update).
func userPlanToSDK(ctx context.Context, plan UserResourceModel, createdTime, updatedTime, internalID string) (*casdoorsdk.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	groups, d := stringListToSDK(ctx, plan.Groups)
	diags.Append(d...)
	address, d := stringListToSDK(ctx, plan.Address)
	diags.Append(d...)
	recoveryCodes, d := stringListToSDK(ctx, plan.RecoveryCodes)
	diags.Append(d...)

	var properties map[string]string
	if !plan.Properties.IsNull() && !plan.Properties.IsUnknown() {
		properties = make(map[string]string)
		diags.Append(plan.Properties.ElementsAs(ctx, &properties, false)...)
	}

	var socialLogins map[string]string
	if !plan.SocialLogins.IsNull() && !plan.SocialLogins.IsUnknown() {
		socialLogins = make(map[string]string)
		diags.Append(plan.SocialLogins.ElementsAs(ctx, &socialLogins, false)...)
	}

	// Extract nested list: addresses
	addresses := make([]*casdoorsdk.Address, 0)
	if !plan.Addresses.IsNull() && !plan.Addresses.IsUnknown() {
		var models []AddressModel
		diags.Append(plan.Addresses.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			addresses = append(addresses, &casdoorsdk.Address{
				Tag:     m.Tag.ValueString(),
				Line1:   m.Line1.ValueString(),
				Line2:   m.Line2.ValueString(),
				City:    m.City.ValueString(),
				State:   m.State.ValueString(),
				ZipCode: m.ZipCode.ValueString(),
				Region:  m.Region.ValueString(),
			})
		}
	}

	// Extract nested list: managed_accounts
	managedAccounts := make([]casdoorsdk.ManagedAccount, 0)
	if !plan.ManagedAccounts.IsNull() && !plan.ManagedAccounts.IsUnknown() {
		var models []ManagedAccountModel
		diags.Append(plan.ManagedAccounts.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			managedAccounts = append(managedAccounts, casdoorsdk.ManagedAccount{
				Application: m.Application.ValueString(),
				Username:    m.Username.ValueString(),
				Password:    m.Password.ValueString(),
				SigninUrl:   m.SigninUrl.ValueString(),
			})
		}
	}

	// Extract nested list: mfa_accounts
	mfaAccounts := make([]casdoorsdk.MfaAccount, 0)
	if !plan.MfaAccounts.IsNull() && !plan.MfaAccounts.IsUnknown() {
		var models []MfaAccountModel
		diags.Append(plan.MfaAccounts.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			mfaAccounts = append(mfaAccounts, casdoorsdk.MfaAccount{
				AccountName: m.AccountName.ValueString(),
				Issuer:      m.Issuer.ValueString(),
				SecretKey:   m.SecretKey.ValueString(),
				Origin:      m.Origin.ValueString(),
			})
		}
	}

	// Extract nested list: mfa_items
	mfaItems := make([]*casdoorsdk.MfaItem, 0)
	if !plan.MfaItems.IsNull() && !plan.MfaItems.IsUnknown() {
		var models []MfaItemModel
		diags.Append(plan.MfaItems.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			mfaItems = append(mfaItems, &casdoorsdk.MfaItem{
				Name: m.Name.ValueString(),
				Rule: m.Rule.ValueString(),
			})
		}
	}

	// Extract nested list: face_ids
	faceIds := make([]*casdoorsdk.FaceId, 0)
	if !plan.FaceIds.IsNull() && !plan.FaceIds.IsUnknown() {
		var models []FaceIdModel
		diags.Append(plan.FaceIds.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			var faceIdData []float64
			if !m.FaceIdData.IsNull() && !m.FaceIdData.IsUnknown() {
				diags.Append(m.FaceIdData.ElementsAs(ctx, &faceIdData, false)...)
			}
			faceIds = append(faceIds, &casdoorsdk.FaceId{
				Name:       m.Name.ValueString(),
				FaceIdData: faceIdData,
				ImageUrl:   m.ImageUrl.ValueString(),
			})
		}
	}

	// Extract nested list: cart
	cart := make([]casdoorsdk.ProductInfo, 0)
	if !plan.Cart.IsNull() && !plan.Cart.IsUnknown() {
		var models []ProductInfoModel
		diags.Append(plan.Cart.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			cart = append(cart, casdoorsdk.ProductInfo{
				Owner:       m.Owner.ValueString(),
				Name:        m.Name.ValueString(),
				DisplayName: m.DisplayName.ValueString(),
				Image:       m.Image.ValueString(),
				Detail:      m.Detail.ValueString(),
				Price:       m.Price.ValueFloat64(),
				Currency:    m.Currency.ValueString(),
				IsRecharge:  m.IsRecharge.ValueBool(),
				Quantity:    int(m.Quantity.ValueInt64()),
				PricingName: m.PricingName.ValueString(),
				PlanName:    m.PlanName.ValueString(),
			})
		}
	}

	if diags.HasError() {
		return nil, diags
	}

	user := &casdoorsdk.User{
		Owner:                plan.Owner.ValueString(),
		Name:                 plan.Name.ValueString(),
		Id:                   internalID,
		CreatedTime:          createdTime,
		UpdatedTime:          updatedTime,
		Type:                 plan.Type.ValueString(),
		Password:             plan.Password.ValueString(),
		PasswordType:         plan.PasswordType.ValueString(),
		DisplayName:          plan.DisplayName.ValueString(),
		FirstName:            plan.FirstName.ValueString(),
		LastName:             plan.LastName.ValueString(),
		Avatar:               plan.Avatar.ValueString(),
		Email:                plan.Email.ValueString(),
		EmailVerified:        plan.EmailVerified.ValueBool(),
		Phone:                plan.Phone.ValueString(),
		CountryCode:          plan.CountryCode.ValueString(),
		Region:               plan.Region.ValueString(),
		Location:             plan.Location.ValueString(),
		Affiliation:          plan.Affiliation.ValueString(),
		Title:                plan.Title.ValueString(),
		Homepage:             plan.Homepage.ValueString(),
		Bio:                  plan.Bio.ValueString(),
		Tag:                  plan.Tag.ValueString(),
		Language:             plan.Language.ValueString(),
		Gender:               plan.Gender.ValueString(),
		Birthday:             plan.Birthday.ValueString(),
		Education:            plan.Education.ValueString(),
		Score:                int(plan.Score.ValueInt64()),
		Karma:                int(plan.Karma.ValueInt64()),
		Ranking:              int(plan.Ranking.ValueInt64()),
		IsAdmin:              plan.IsAdmin.ValueBool(),
		IsForbidden:          plan.IsForbidden.ValueBool(),
		IsDeleted:            plan.IsDeleted.ValueBool(),
		SignupApplication:    plan.SignupApplication.ValueString(),
		ExternalId:           plan.ExternalId.ValueString(),
		PasswordSalt:         plan.PasswordSalt.ValueString(),
		AvatarType:           plan.AvatarType.ValueString(),
		PermanentAvatar:      plan.PermanentAvatar.ValueString(),
		Address:              address,
		Addresses:            addresses,
		IdCardType:           plan.IdCardType.ValueString(),
		IdCard:               plan.IdCard.ValueString(),
		RealName:             plan.RealName.ValueString(),
		IsVerified:           plan.IsVerified.ValueBool(),
		Balance:              plan.Balance.ValueFloat64(),
		BalanceCredit:        plan.BalanceCredit.ValueFloat64(),
		Currency:             plan.Currency.ValueString(),
		BalanceCurrency:      plan.BalanceCurrency.ValueString(),
		RegisterType:         plan.RegisterType.ValueString(),
		RegisterSource:       plan.RegisterSource.ValueString(),
		AccessKey:            plan.AccessKey.ValueString(),
		AccessSecret:         plan.AccessSecret.ValueString(),
		AccessToken:          plan.AccessToken.ValueString(),
		OriginalToken:        plan.OriginalToken.ValueString(),
		OriginalRefreshToken: plan.OriginalRefreshToken.ValueString(),
		Invitation:           plan.Invitation.ValueString(),
		InvitationCode:       plan.InvitationCode.ValueString(),
		Ldap:                 plan.Ldap.ValueString(),
		Properties:           properties,
		NeedUpdatePassword:   plan.NeedUpdatePassword.ValueBool(),
		PreferredMfaType:     plan.PreferredMfaType.ValueString(),
		RecoveryCodes:        recoveryCodes,
		TotpSecret:           plan.TotpSecret.ValueString(),
		MfaPhoneEnabled:      plan.MfaPhoneEnabled.ValueBool(),
		MfaEmailEnabled:      plan.MfaEmailEnabled.ValueBool(),
		MfaRadiusEnabled:     plan.MfaRadiusEnabled.ValueBool(),
		MfaRadiusUsername:    plan.MfaRadiusUsername.ValueString(),
		MfaRadiusProvider:    plan.MfaRadiusProvider.ValueString(),
		MfaPushEnabled:       plan.MfaPushEnabled.ValueBool(),
		MfaPushReceiver:      plan.MfaPushReceiver.ValueString(),
		MfaPushProvider:      plan.MfaPushProvider.ValueString(),
		MfaRememberDeadline:  plan.MfaRememberDeadline.ValueString(),
		IpWhitelist:          plan.IpWhitelist.ValueString(),
		ManagedAccounts:      managedAccounts,
		MfaAccounts:          mfaAccounts,
		MfaItems:             mfaItems,
		FaceIds:              faceIds,
		Cart:                 cart,
		Groups:               groups,
	}

	// Set social login fields on the SDK struct.
	for _, f := range socialLoginFields {
		if v, ok := socialLogins[f.Key]; ok {
			f.Set(user, v)
		}
	}

	return user, diags
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdTime := plan.CreatedTime.ValueString()
	if createdTime == "" {
		createdTime = time.Now().UTC().Format(time.RFC3339)
	}

	user, diags := userPlanToSDK(ctx, plan, createdTime, createdTime, "")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ok, err := r.client.AddUser(user)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("creating user %q", plan.Name.ValueString())) {
		return
	}

	// Read back the user to get the generated ID.
	createdUser, err := r.client.GetUser(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User After Create",
			fmt.Sprintf("Could not read user %q after creation: %s", plan.Name.ValueString(), err),
		)
		return
	}

	if createdUser != nil {
		plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
		plan.CreatedTime = types.StringValue(createdUser.CreatedTime)
		plan.UpdatedTime = types.StringValue(createdUser.UpdatedTime)
		plan.DeletedTime = types.StringValue(createdUser.DeletedTime)
		plan.IsDefaultAvatar = types.BoolValue(createdUser.IsDefaultAvatar)
		plan.IsOnline = types.BoolValue(createdUser.IsOnline)
		plan.Hash = types.StringValue(createdUser.Hash)
		plan.PreHash = types.StringValue(createdUser.PreHash)
		plan.CreatedIp = types.StringValue(createdUser.CreatedIp)
		plan.LastSigninTime = types.StringValue(createdUser.LastSigninTime)
		plan.LastSigninIp = types.StringValue(createdUser.LastSigninIp)
		plan.LastChangePasswordTime = types.StringValue(createdUser.LastChangePasswordTime)
		plan.LastSigninWrongTime = types.StringValue(createdUser.LastSigninWrongTime)
		plan.SigninWrongTimes = types.Int64Value(int64(createdUser.SigninWrongTimes))
		plan.BalanceCurrency = types.StringValue(createdUser.BalanceCurrency)
		plan.RegisterType = types.StringValue(createdUser.RegisterType)
		// Preserve server-generated sensitive values.
		if createdUser.PasswordSalt != "***" {
			plan.PasswordSalt = types.StringValue(createdUser.PasswordSalt)
		}
		if createdUser.AccessKey != "***" {
			plan.AccessKey = types.StringValue(createdUser.AccessKey)
		}
		if createdUser.AccessSecret != "***" {
			plan.AccessSecret = types.StringValue(createdUser.AccessSecret)
		}
		if createdUser.TotpSecret != "***" {
			plan.TotpSecret = types.StringValue(createdUser.TotpSecret)
		}
		if createdUser.AccessToken != "***" {
			plan.AccessToken = types.StringValue(createdUser.AccessToken)
		}
		if createdUser.OriginalToken != "***" {
			plan.OriginalToken = types.StringValue(createdUser.OriginalToken)
		}
		if createdUser.OriginalRefreshToken != "***" {
			plan.OriginalRefreshToken = types.StringValue(createdUser.OriginalRefreshToken)
		}
	} else {
		// GetUser uses the provider's OrganizationName which may differ from the
		// user's owner. Set computed fields to defaults to avoid unknown values.
		plan.ID = types.StringValue(plan.Owner.ValueString() + "/" + plan.Name.ValueString())
		plan.CreatedTime = types.StringValue(createdTime)
		plan.UpdatedTime = types.StringValue(createdTime)
		plan.DeletedTime = types.StringValue("")
		plan.IsDefaultAvatar = types.BoolValue(false)
		plan.IsOnline = types.BoolValue(false)
		plan.Hash = types.StringValue("")
		plan.PreHash = types.StringValue("")
		plan.CreatedIp = types.StringValue("")
		plan.LastSigninTime = types.StringValue("")
		plan.LastSigninIp = types.StringValue("")
		plan.LastChangePasswordTime = types.StringValue("")
		plan.LastSigninWrongTime = types.StringValue("")
		plan.SigninWrongTimes = types.Int64Value(0)
		plan.BalanceCurrency = types.StringValue("")
		plan.RegisterType = types.StringValue("")
	}

	plan.Groups, diags = stringListFromSDK(ctx, user.Groups)
	resp.Diagnostics.Append(diags...)
	plan.Address, diags = stringListFromSDK(ctx, user.Address)
	resp.Diagnostics.Append(diags...)
	plan.RecoveryCodes, diags = stringListFromSDK(ctx, user.RecoveryCodes)
	resp.Diagnostics.Append(diags...)
	if len(user.Properties) == 0 {
		plan.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}
	if len(user.Addresses) == 0 {
		plan.Addresses = types.ListNull(types.ObjectType{AttrTypes: AddressAttrTypes()})
	}
	if len(user.ManagedAccounts) == 0 {
		plan.ManagedAccounts = types.ListNull(types.ObjectType{AttrTypes: ManagedAccountAttrTypes()})
	}
	if len(user.MfaAccounts) == 0 {
		plan.MfaAccounts = types.ListNull(types.ObjectType{AttrTypes: MfaAccountAttrTypes()})
	}
	if len(user.MfaItems) == 0 {
		plan.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}
	if len(user.FaceIds) == 0 {
		plan.FaceIds = types.ListNull(types.ObjectType{AttrTypes: FaceIdAttrTypes()})
	}
	if len(user.Cart) == 0 {
		plan.Cart = types.ListNull(types.ObjectType{AttrTypes: ProductInfoAttrTypes()})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := getByOwnerName[casdoorsdk.User](r.client, "get-user", state.Owner.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User",
			fmt.Sprintf("Could not read user %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Owner = types.StringValue(user.Owner)
	state.Name = types.StringValue(user.Name)
	state.ID = types.StringValue(user.Owner + "/" + user.Name)
	state.Type = types.StringValue(user.Type)
	// Password is write-only; never read it back from the server.
	state.PasswordType = types.StringValue(user.PasswordType)
	state.DisplayName = types.StringValue(user.DisplayName)
	state.FirstName = types.StringValue(user.FirstName)
	state.LastName = types.StringValue(user.LastName)
	state.Avatar = types.StringValue(user.Avatar)
	state.Email = types.StringValue(user.Email)
	state.EmailVerified = types.BoolValue(user.EmailVerified)
	state.Phone = types.StringValue(user.Phone)
	state.CountryCode = types.StringValue(user.CountryCode)
	state.Region = types.StringValue(user.Region)
	state.Location = types.StringValue(user.Location)
	state.Affiliation = types.StringValue(user.Affiliation)
	state.Title = types.StringValue(user.Title)
	state.Homepage = types.StringValue(user.Homepage)
	state.Bio = types.StringValue(user.Bio)
	state.Tag = types.StringValue(user.Tag)
	state.Language = types.StringValue(user.Language)
	state.Gender = types.StringValue(user.Gender)
	state.Birthday = types.StringValue(user.Birthday)
	state.Education = types.StringValue(user.Education)
	state.Score = types.Int64Value(int64(user.Score))
	state.Karma = types.Int64Value(int64(user.Karma))
	state.Ranking = types.Int64Value(int64(user.Ranking))
	state.IsAdmin = types.BoolValue(user.IsAdmin)
	state.IsForbidden = types.BoolValue(user.IsForbidden)
	state.IsDeleted = types.BoolValue(user.IsDeleted)
	state.SignupApplication = types.StringValue(user.SignupApplication)
	state.CreatedTime = types.StringValue(user.CreatedTime)
	state.UpdatedTime = types.StringValue(user.UpdatedTime)
	state.DeletedTime = types.StringValue(user.DeletedTime)
	state.ExternalId = types.StringValue(user.ExternalId)
	if user.PasswordSalt != "***" {
		state.PasswordSalt = types.StringValue(user.PasswordSalt)
	}
	state.AvatarType = types.StringValue(user.AvatarType)
	state.PermanentAvatar = types.StringValue(user.PermanentAvatar)
	state.IdCardType = types.StringValue(user.IdCardType)
	state.IdCard = types.StringValue(user.IdCard)
	state.RealName = types.StringValue(user.RealName)
	state.IsVerified = types.BoolValue(user.IsVerified)
	state.IsDefaultAvatar = types.BoolValue(user.IsDefaultAvatar)
	state.IsOnline = types.BoolValue(user.IsOnline)
	state.Hash = types.StringValue(user.Hash)
	state.PreHash = types.StringValue(user.PreHash)
	state.Balance = types.Float64Value(user.Balance)
	state.BalanceCredit = types.Float64Value(user.BalanceCredit)
	state.Currency = types.StringValue(user.Currency)
	state.BalanceCurrency = types.StringValue(user.BalanceCurrency)
	state.RegisterType = types.StringValue(user.RegisterType)
	state.RegisterSource = types.StringValue(user.RegisterSource)
	if user.AccessKey != "***" {
		state.AccessKey = types.StringValue(user.AccessKey)
	}
	if user.AccessSecret != "***" {
		state.AccessSecret = types.StringValue(user.AccessSecret)
	}
	if user.AccessToken != "***" {
		state.AccessToken = types.StringValue(user.AccessToken)
	}
	if user.OriginalToken != "***" {
		state.OriginalToken = types.StringValue(user.OriginalToken)
	}
	if user.OriginalRefreshToken != "***" {
		state.OriginalRefreshToken = types.StringValue(user.OriginalRefreshToken)
	}
	state.CreatedIp = types.StringValue(user.CreatedIp)
	state.LastSigninTime = types.StringValue(user.LastSigninTime)
	state.LastSigninIp = types.StringValue(user.LastSigninIp)
	state.Invitation = types.StringValue(user.Invitation)
	state.InvitationCode = types.StringValue(user.InvitationCode)
	state.Ldap = types.StringValue(user.Ldap)
	state.NeedUpdatePassword = types.BoolValue(user.NeedUpdatePassword)
	state.LastChangePasswordTime = types.StringValue(user.LastChangePasswordTime)
	state.LastSigninWrongTime = types.StringValue(user.LastSigninWrongTime)
	state.SigninWrongTimes = types.Int64Value(int64(user.SigninWrongTimes))
	state.PreferredMfaType = types.StringValue(user.PreferredMfaType)
	if user.TotpSecret != "***" {
		state.TotpSecret = types.StringValue(user.TotpSecret)
	}
	state.MfaPhoneEnabled = types.BoolValue(user.MfaPhoneEnabled)
	state.MfaEmailEnabled = types.BoolValue(user.MfaEmailEnabled)
	state.MfaRadiusEnabled = types.BoolValue(user.MfaRadiusEnabled)
	state.MfaRadiusUsername = types.StringValue(user.MfaRadiusUsername)
	state.MfaRadiusProvider = types.StringValue(user.MfaRadiusProvider)
	state.MfaPushEnabled = types.BoolValue(user.MfaPushEnabled)
	state.MfaPushReceiver = types.StringValue(user.MfaPushReceiver)
	state.MfaPushProvider = types.StringValue(user.MfaPushProvider)
	state.MfaRememberDeadline = types.StringValue(user.MfaRememberDeadline)
	state.IpWhitelist = types.StringValue(user.IpWhitelist)

	// Read social logins map.
	socialLoginsMap := make(map[string]attr.Value)
	for _, f := range socialLoginFields {
		v := f.Get(user)
		if v != "" {
			socialLoginsMap[f.Key] = types.StringValue(v)
		}
	}
	if len(socialLoginsMap) > 0 {
		state.SocialLogins = types.MapValueMust(types.StringType, socialLoginsMap)
	} else {
		state.SocialLogins = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	// Read address ([]string).
	var diags diag.Diagnostics
	state.Address, diags = stringListFromSDK(ctx, user.Address)
	resp.Diagnostics.Append(diags...)

	// Read addresses ([]*Address).
	if len(user.Addresses) > 0 {
		objList := make([]attr.Value, 0, len(user.Addresses))
		for _, item := range user.Addresses {
			obj, diags := types.ObjectValue(AddressAttrTypes(), map[string]attr.Value{
				"tag":      types.StringValue(item.Tag),
				"line1":    types.StringValue(item.Line1),
				"line2":    types.StringValue(item.Line2),
				"city":     types.StringValue(item.City),
				"state":    types.StringValue(item.State),
				"zip_code": types.StringValue(item.ZipCode),
				"region":   types.StringValue(item.Region),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: AddressAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.Addresses = list
	} else {
		state.Addresses = types.ListNull(types.ObjectType{AttrTypes: AddressAttrTypes()})
	}

	// Read managed_accounts.
	if len(user.ManagedAccounts) > 0 {
		objList := make([]attr.Value, 0, len(user.ManagedAccounts))
		for _, item := range user.ManagedAccounts {
			obj, diags := types.ObjectValue(ManagedAccountAttrTypes(), map[string]attr.Value{
				"application": types.StringValue(item.Application),
				"username":    types.StringValue(item.Username),
				"password":    types.StringValue(item.Password),
				"signin_url":  types.StringValue(item.SigninUrl),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: ManagedAccountAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.ManagedAccounts = list
	} else {
		state.ManagedAccounts = types.ListNull(types.ObjectType{AttrTypes: ManagedAccountAttrTypes()})
	}

	// Read mfa_accounts.
	if len(user.MfaAccounts) > 0 {
		objList := make([]attr.Value, 0, len(user.MfaAccounts))
		for _, item := range user.MfaAccounts {
			obj, diags := types.ObjectValue(MfaAccountAttrTypes(), map[string]attr.Value{
				"account_name": types.StringValue(item.AccountName),
				"issuer":       types.StringValue(item.Issuer),
				"secret_key":   types.StringValue(item.SecretKey),
				"origin":       types.StringValue(item.Origin),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: MfaAccountAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.MfaAccounts = list
	} else {
		state.MfaAccounts = types.ListNull(types.ObjectType{AttrTypes: MfaAccountAttrTypes()})
	}

	// Read mfa_items.
	if len(user.MfaItems) > 0 {
		objList := make([]attr.Value, 0, len(user.MfaItems))
		for _, item := range user.MfaItems {
			obj, diags := types.ObjectValue(MfaItemAttrTypes(), map[string]attr.Value{
				"name": types.StringValue(item.Name),
				"rule": types.StringValue(item.Rule),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: MfaItemAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.MfaItems = list
	} else {
		state.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}

	// Read face_ids.
	if len(user.FaceIds) > 0 {
		objList := make([]attr.Value, 0, len(user.FaceIds))
		for _, item := range user.FaceIds {
			var faceIdDataVal attr.Value
			if len(item.FaceIdData) > 0 {
				floatVals := make([]attr.Value, 0, len(item.FaceIdData))
				for _, f := range item.FaceIdData {
					floatVals = append(floatVals, types.Float64Value(f))
				}
				faceIdDataVal = types.ListValueMust(types.Float64Type, floatVals)
			} else {
				faceIdDataVal = types.ListNull(types.Float64Type)
			}
			obj, diags := types.ObjectValue(FaceIdAttrTypes(), map[string]attr.Value{
				"name":         types.StringValue(item.Name),
				"face_id_data": faceIdDataVal,
				"image_url":    types.StringValue(item.ImageUrl),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: FaceIdAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.FaceIds = list
	} else {
		state.FaceIds = types.ListNull(types.ObjectType{AttrTypes: FaceIdAttrTypes()})
	}

	// Read cart.
	if len(user.Cart) > 0 {
		objList := make([]attr.Value, 0, len(user.Cart))
		for _, item := range user.Cart {
			obj, diags := types.ObjectValue(ProductInfoAttrTypes(), map[string]attr.Value{
				"owner":        types.StringValue(item.Owner),
				"name":         types.StringValue(item.Name),
				"display_name": types.StringValue(item.DisplayName),
				"image":        types.StringValue(item.Image),
				"detail":       types.StringValue(item.Detail),
				"price":        types.Float64Value(item.Price),
				"currency":     types.StringValue(item.Currency),
				"is_recharge":  types.BoolValue(item.IsRecharge),
				"quantity":     types.Int64Value(int64(item.Quantity)),
				"pricing_name": types.StringValue(item.PricingName),
				"plan_name":    types.StringValue(item.PlanName),
			})
			resp.Diagnostics.Append(diags...)
			objList = append(objList, obj)
		}
		list, diags := types.ListValue(types.ObjectType{AttrTypes: ProductInfoAttrTypes()}, objList)
		resp.Diagnostics.Append(diags...)
		state.Cart = list
	} else {
		state.Cart = types.ListNull(types.ObjectType{AttrTypes: ProductInfoAttrTypes()})
	}

	// Read recovery_codes.
	state.RecoveryCodes, diags = stringListFromSDK(ctx, user.RecoveryCodes)
	resp.Diagnostics.Append(diags...)

	// Read properties.
	if len(user.Properties) > 0 {
		props, diags := types.MapValueFrom(ctx, types.StringType, user.Properties)
		resp.Diagnostics.Append(diags...)
		state.Properties = props
	} else {
		state.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	// Read groups.
	state.Groups, diags = stringListFromSDK(ctx, user.Groups)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups := make([]string, 0)
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	}
	address := make([]string, 0)
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		resp.Diagnostics.Append(plan.Address.ElementsAs(ctx, &address, false)...)
	}
	recoveryCodes := make([]string, 0)
	if !plan.RecoveryCodes.IsNull() && !plan.RecoveryCodes.IsUnknown() {
		resp.Diagnostics.Append(plan.RecoveryCodes.ElementsAs(ctx, &recoveryCodes, false)...)
	}
	var properties map[string]string
	if !plan.Properties.IsNull() && !plan.Properties.IsUnknown() {
		properties = make(map[string]string)
		resp.Diagnostics.Append(plan.Properties.ElementsAs(ctx, &properties, false)...)
	}
	var socialLogins map[string]string
	if !plan.SocialLogins.IsNull() && !plan.SocialLogins.IsUnknown() {
		socialLogins = make(map[string]string)
		resp.Diagnostics.Append(plan.SocialLogins.ElementsAs(ctx, &socialLogins, false)...)
	}

	// Extract nested list: addresses
	addresses := make([]*casdoorsdk.Address, 0)
	if !plan.Addresses.IsNull() && !plan.Addresses.IsUnknown() {
		var models []AddressModel
		resp.Diagnostics.Append(plan.Addresses.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			addresses = append(addresses, &casdoorsdk.Address{
				Tag:     m.Tag.ValueString(),
				Line1:   m.Line1.ValueString(),
				Line2:   m.Line2.ValueString(),
				City:    m.City.ValueString(),
				State:   m.State.ValueString(),
				ZipCode: m.ZipCode.ValueString(),
				Region:  m.Region.ValueString(),
			})
		}
	}

	// Extract nested list: managed_accounts
	managedAccounts := make([]casdoorsdk.ManagedAccount, 0)
	if !plan.ManagedAccounts.IsNull() && !plan.ManagedAccounts.IsUnknown() {
		var models []ManagedAccountModel
		resp.Diagnostics.Append(plan.ManagedAccounts.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			managedAccounts = append(managedAccounts, casdoorsdk.ManagedAccount{
				Application: m.Application.ValueString(),
				Username:    m.Username.ValueString(),
				Password:    m.Password.ValueString(),
				SigninUrl:   m.SigninUrl.ValueString(),
			})
		}
	}

	// Extract nested list: mfa_accounts
	mfaAccounts := make([]casdoorsdk.MfaAccount, 0)
	if !plan.MfaAccounts.IsNull() && !plan.MfaAccounts.IsUnknown() {
		var models []MfaAccountModel
		resp.Diagnostics.Append(plan.MfaAccounts.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			mfaAccounts = append(mfaAccounts, casdoorsdk.MfaAccount{
				AccountName: m.AccountName.ValueString(),
				Issuer:      m.Issuer.ValueString(),
				SecretKey:   m.SecretKey.ValueString(),
				Origin:      m.Origin.ValueString(),
			})
		}
	}

	// Extract nested list: mfa_items
	mfaItems := make([]*casdoorsdk.MfaItem, 0)
	if !plan.MfaItems.IsNull() && !plan.MfaItems.IsUnknown() {
		var models []MfaItemModel
		resp.Diagnostics.Append(plan.MfaItems.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			mfaItems = append(mfaItems, &casdoorsdk.MfaItem{
				Name: m.Name.ValueString(),
				Rule: m.Rule.ValueString(),
			})
		}
	}

	// Extract nested list: face_ids
	faceIds := make([]*casdoorsdk.FaceId, 0)
	if !plan.FaceIds.IsNull() && !plan.FaceIds.IsUnknown() {
		var models []FaceIdModel
		resp.Diagnostics.Append(plan.FaceIds.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			var faceIdData []float64
			if !m.FaceIdData.IsNull() && !m.FaceIdData.IsUnknown() {
				resp.Diagnostics.Append(m.FaceIdData.ElementsAs(ctx, &faceIdData, false)...)
			}
			faceIds = append(faceIds, &casdoorsdk.FaceId{
				Name:       m.Name.ValueString(),
				FaceIdData: faceIdData,
				ImageUrl:   m.ImageUrl.ValueString(),
			})
		}
	}

	// Extract nested list: cart
	cart := make([]casdoorsdk.ProductInfo, 0)
	if !plan.Cart.IsNull() && !plan.Cart.IsUnknown() {
		var models []ProductInfoModel
		resp.Diagnostics.Append(plan.Cart.ElementsAs(ctx, &models, false)...)
		for _, m := range models {
			cart = append(cart, casdoorsdk.ProductInfo{
				Owner:       m.Owner.ValueString(),
				Name:        m.Name.ValueString(),
				DisplayName: m.DisplayName.ValueString(),
				Image:       m.Image.ValueString(),
				Detail:      m.Detail.ValueString(),
				Price:       m.Price.ValueFloat64(),
				Currency:    m.Currency.ValueString(),
				IsRecharge:  m.IsRecharge.ValueBool(),
				Quantity:    int(m.Quantity.ValueInt64()),
				PricingName: m.PricingName.ValueString(),
				PlanName:    m.PlanName.ValueString(),
			})
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Read the existing user to preserve the Casdoor-internal Id field,
	// which is immutable and must not change during updates.
	existingUser, err := r.client.GetUser(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User Before Update",
			fmt.Sprintf("Could not read user %q before update: %s", plan.Name.ValueString(), err),
		)
		return
	}

	var internalID string
	if existingUser != nil {
		internalID = existingUser.Id
	}

	user := &casdoorsdk.User{
		Owner:                plan.Owner.ValueString(),
		Name:                 plan.Name.ValueString(),
		Id:                   internalID,
		Type:                 plan.Type.ValueString(),
		Password:             plan.Password.ValueString(),
		PasswordType:         plan.PasswordType.ValueString(),
		DisplayName:          plan.DisplayName.ValueString(),
		FirstName:            plan.FirstName.ValueString(),
		LastName:             plan.LastName.ValueString(),
		Avatar:               plan.Avatar.ValueString(),
		Email:                plan.Email.ValueString(),
		EmailVerified:        plan.EmailVerified.ValueBool(),
		Phone:                plan.Phone.ValueString(),
		CountryCode:          plan.CountryCode.ValueString(),
		Region:               plan.Region.ValueString(),
		Location:             plan.Location.ValueString(),
		Affiliation:          plan.Affiliation.ValueString(),
		Title:                plan.Title.ValueString(),
		Homepage:             plan.Homepage.ValueString(),
		Bio:                  plan.Bio.ValueString(),
		Tag:                  plan.Tag.ValueString(),
		Language:             plan.Language.ValueString(),
		Gender:               plan.Gender.ValueString(),
		Birthday:             plan.Birthday.ValueString(),
		Education:            plan.Education.ValueString(),
		Score:                int(plan.Score.ValueInt64()),
		Karma:                int(plan.Karma.ValueInt64()),
		Ranking:              int(plan.Ranking.ValueInt64()),
		IsAdmin:              plan.IsAdmin.ValueBool(),
		IsForbidden:          plan.IsForbidden.ValueBool(),
		IsDeleted:            plan.IsDeleted.ValueBool(),
		SignupApplication:    plan.SignupApplication.ValueString(),
		CreatedTime:          plan.CreatedTime.ValueString(),
		ExternalId:           plan.ExternalId.ValueString(),
		PasswordSalt:         plan.PasswordSalt.ValueString(),
		AvatarType:           plan.AvatarType.ValueString(),
		PermanentAvatar:      plan.PermanentAvatar.ValueString(),
		Address:              address,
		Addresses:            addresses,
		IdCardType:           plan.IdCardType.ValueString(),
		IdCard:               plan.IdCard.ValueString(),
		RealName:             plan.RealName.ValueString(),
		IsVerified:           plan.IsVerified.ValueBool(),
		Balance:              plan.Balance.ValueFloat64(),
		BalanceCredit:        plan.BalanceCredit.ValueFloat64(),
		Currency:             plan.Currency.ValueString(),
		BalanceCurrency:      plan.BalanceCurrency.ValueString(),
		RegisterType:         plan.RegisterType.ValueString(),
		RegisterSource:       plan.RegisterSource.ValueString(),
		AccessKey:            plan.AccessKey.ValueString(),
		AccessSecret:         plan.AccessSecret.ValueString(),
		AccessToken:          plan.AccessToken.ValueString(),
		OriginalToken:        plan.OriginalToken.ValueString(),
		OriginalRefreshToken: plan.OriginalRefreshToken.ValueString(),
		Invitation:           plan.Invitation.ValueString(),
		InvitationCode:       plan.InvitationCode.ValueString(),
		Ldap:                 plan.Ldap.ValueString(),
		Properties:           properties,
		NeedUpdatePassword:   plan.NeedUpdatePassword.ValueBool(),
		PreferredMfaType:     plan.PreferredMfaType.ValueString(),
		RecoveryCodes:        recoveryCodes,
		TotpSecret:           plan.TotpSecret.ValueString(),
		MfaPhoneEnabled:      plan.MfaPhoneEnabled.ValueBool(),
		MfaEmailEnabled:      plan.MfaEmailEnabled.ValueBool(),
		MfaRadiusEnabled:     plan.MfaRadiusEnabled.ValueBool(),
		MfaRadiusUsername:    plan.MfaRadiusUsername.ValueString(),
		MfaRadiusProvider:    plan.MfaRadiusProvider.ValueString(),
		MfaPushEnabled:       plan.MfaPushEnabled.ValueBool(),
		MfaPushReceiver:      plan.MfaPushReceiver.ValueString(),
		MfaPushProvider:      plan.MfaPushProvider.ValueString(),
		MfaRememberDeadline:  plan.MfaRememberDeadline.ValueString(),
		IpWhitelist:          plan.IpWhitelist.ValueString(),
		ManagedAccounts:      managedAccounts,
		MfaAccounts:          mfaAccounts,
		MfaItems:             mfaItems,
		FaceIds:              faceIds,
		Cart:                 cart,
		Groups:               groups,
	}

	// Set social login fields on the SDK struct.
	for _, f := range socialLoginFields {
		if v, ok := socialLogins[f.Key]; ok {
			f.Set(user, v)
		}
	}

	ok, err := r.client.UpdateUser(user)
	if sdkError(&resp.Diagnostics, ok, err, fmt.Sprintf("updating user %q", plan.Name.ValueString())) {
		return
	}

	if len(groups) == 0 {
		plan.Groups = types.ListNull(types.StringType)
	}
	if len(address) == 0 {
		plan.Address = types.ListNull(types.StringType)
	}
	if len(recoveryCodes) == 0 {
		plan.RecoveryCodes = types.ListNull(types.StringType)
	}
	if len(properties) == 0 {
		plan.Properties = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}
	if len(addresses) == 0 {
		plan.Addresses = types.ListNull(types.ObjectType{AttrTypes: AddressAttrTypes()})
	}
	if len(managedAccounts) == 0 {
		plan.ManagedAccounts = types.ListNull(types.ObjectType{AttrTypes: ManagedAccountAttrTypes()})
	}
	if len(mfaAccounts) == 0 {
		plan.MfaAccounts = types.ListNull(types.ObjectType{AttrTypes: MfaAccountAttrTypes()})
	}
	if len(mfaItems) == 0 {
		plan.MfaItems = types.ListNull(types.ObjectType{AttrTypes: MfaItemAttrTypes()})
	}
	if len(faceIds) == 0 {
		plan.FaceIds = types.ListNull(types.ObjectType{AttrTypes: FaceIdAttrTypes()})
	}
	if len(cart) == 0 {
		plan.Cart = types.ListNull(types.ObjectType{AttrTypes: ProductInfoAttrTypes()})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &casdoorsdk.User{
		Owner: state.Owner.ValueString(),
		Name:  state.Name.ValueString(),
	}

	_, err := r.client.DeleteUser(user)
	if err != nil {
		// Casdoor returns "session is nil" when deleting users from the built-in
		// organization. This is a known server-side bug; treat it as a warning
		// rather than a hard error.
		if strings.Contains(err.Error(), "session is nil") {
			resp.Diagnostics.AddWarning(
				"User Deletion Warning",
				fmt.Sprintf("Casdoor returned 'session is nil' when deleting user %q. "+
					"This is a known Casdoor issue with the built-in organization.", state.Name.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting User",
			fmt.Sprintf("Could not delete user %q: %s", state.Name.ValueString(), err),
		)
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateOwnerName(ctx, req, resp)
}
