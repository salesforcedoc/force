package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ForceCLI/force/config"
	. "github.com/ForceCLI/force/error"
	. "github.com/ForceCLI/force/lib"
)

var cmdExport = &Command{
	Run:   runExport,
	Usage: "export [options] [dir]",
	Short: "Export metadata to a local directory",
	Long: `
Export metadata to a local directory

Export Options
  -w, -warnings  # Display warnings about metadata that cannot be retrieved
  -x, -exclude   # Exclude given metadata type
  -i, -include   # Include given metadata type
  -p, -package   # Include managed packages

Examples:

  force export

  force export org/schema

  force export -x ApexClass -x CustomObject
`,
	MaxExpectedArgs: 1,
}

type metadataList []string

func (i *metadataList) String() string {
	return fmt.Sprint(*i)
}

func (i *metadataList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	showWarnings           bool
	includeManagedPackages bool
	excludeMetadataNames   metadataList
	includeMetadataNames   metadataList
)

func init() {
	cmdExport.Flag.BoolVar(&showWarnings, "w", false, "show warnings")
	cmdExport.Flag.BoolVar(&showWarnings, "warnings", false, "show warnings")
	cmdExport.Flag.Var(&excludeMetadataNames, "x", "exclude metadata type")
	cmdExport.Flag.Var(&excludeMetadataNames, "exclude", "exclude metadata type")
	cmdExport.Flag.Var(&includeMetadataNames, "i", "include only metadata type")
	cmdExport.Flag.Var(&includeMetadataNames, "include", "include only metadata type")
	cmdExport.Flag.BoolVar(&includeManagedPackages, "p", false, "include managed packages")
	cmdExport.Flag.BoolVar(&includeManagedPackages, "package", false, "include managed packages")
}

func runExport(cmd *Command, args []string) {
	// Get path from args if available
	var err error
	var root string
	if len(args) == 1 {
		root, err = filepath.Abs(args[0])
	}
	if err != nil {
		fmt.Printf("Error obtaining file path\n")
		ErrorAndExit(err.Error())
	}
	force, _ := ActiveForce()
	sobjects, err := force.ListSobjects()
	if err != nil {
		ErrorAndExit(err.Error())
	}
	query := make(ForceMetadataQuery, 0)
	customObject := "CustomObject"

	sort.Strings(excludeMetadataNames)
	sort.Strings(includeMetadataNames)

	if !isExcluded(customObject) || isIncluded(customObject) {
		stdObjects := make([]string, 1, len(sobjects)+1)
		stdObjects[0] = "*"
		for _, sobject := range sobjects {
			name := sobject["name"].(string)
			include := true
			if strings.Count(name, "__") > 1 {
				if !includeManagedPackages {
					include = false
				}
			}
			if include && !strings.HasSuffix(name, "Tag") && !strings.HasSuffix(name, "History") &&
				!strings.HasSuffix(name, "Share") && !strings.HasSuffix(name, "ChangeEvent") &&
				!strings.HasSuffix(name, "Feed") {
				stdObjects = append(stdObjects, name)
			}
		}
		stdObjects = append(stdObjects, "Activity")

		query = append(query, ForceMetadataQueryElement{Name: []string{customObject}, Members: stdObjects})
	}

	metadataNames := []string{
		"AccessControlPolicy",
		"AccountForecastSettings",
		"AccountInsightsSettings",
		"AccountIntelligenceSettings",
		"AccountRelationshipShareRule",
		"AccountSettings",
		"AcctMgrTargetSettings",
		"ActionLinkGroupTemplate",
		"ActionPlanTemplate",
		"ActionsSettings",
		"ActivitiesSettings",
		"AddressSettings",
		"AIReplyRecommendationsSettings",
		"AnalyticSnapshot",
		"AnalyticsSettings",
		"AnimationRule",
		"ApexClass",
		"ApexComponent",
		"ApexEmailNotifications",
		"ApexPage",
		"ApexSettings",
		"ApexTestSuite",
		"ApexTrigger",
		"AppAnalyticsSettings",
		"AppExperienceSettings",
		"ApplicationRecordTypeConfig",
		"AppMenu",
		"AppointmentSchedulingPolicy",
		"ApprovalProcess",
		"ArchiveSettings",
		"AssignmentRules",
		"AssistantContextItem",
		"AssistantDefinition",
		"AssistantSkillQuickAction",
		"AssistantSkillSobjectAction",
		"AssistantVersion",
		"Audience",
		"AuraDefinitionBundle",
		"AuthProvider",
		"AutomatedContactsSettings",
		"AutoResponseRules",
		"BatchCalcJobDefinition",
		"BatchProcessJobDefinition",
		"BlacklistedConsumer",
		"BlockchainSettings",
		"Bot",
		"BotSettings",
		"BotVersion",
		"BrandingSet",
		"BusinessHoursSettings",
		"BusinessProcess",
		"BusinessProcessGroup",
		"CallCenter",
		"CallCoachingMediaProvider",
		"CampaignInfluenceModel",
		"CampaignSettings",
		"CanvasMetadata",
		"CareProviderSearchConfig",
		"CareRequestConfiguration",
		"CareSystemFieldMapping",
		"CaseClassificationSettings",
		"CaseSettings",
		"CaseSubjectParticle",
		"Certificate",
		"ChannelLayout",
		"ChannelObjectLinkingRule",
		"ChatterAnswersSettings",
		"ChatterEmailsMDSettings",
		"ChatterExtension",
		"ChatterSettings",
		"CleanDataService",
		"CMSConnectSource",
		"CommandAction",
		"CommunitiesSettings",
		"Community",
		"CommunityTemplateDefinition",
		"CommunityThemeDefinition",
		"CompactLayout",
		"CompanySettings",
		"ConnectedApp",
		"ConnectedAppSettings",
		"ContentAsset",
		"ContentSettings",
		"ContractSettings",
		"ConversationalIntelligenceSettings",
		"CorsWhitelistOrigin",
		"CspTrustedSite",
		"CurrencySettings",
		"CustomApplication",
		"CustomApplicationComponent",
		"CustomerDataPlatformSettings",
		"CustomFeedFilter",
		"CustomField",
		"CustomHelpMenuSection",
		"CustomLabels",
		"CustomMetadata",
		"CustomNotificationType",
		"CustomObjectTranslation",
		"CustomPageWebLink",
		"CustomPermission",
		"CustomSite",
		"CustomTab",
		"DashboardFolder",
		"DataCategoryGroup",
		"DataDotComSettings",
		"DataSourceObject",
		"DecisionTable",
		"DecisionTableDatasetLink",
		"DelegateGroup",
		"DeploymentSettings",
		"DevHubSettings",
		"DiscoverySettings",
		"DocumentChecklistSettings",
		"DocumentFolder",
		"DocumentType",
		"DuplicateRule",
		"DynamicTrigger",
		"EACSettings",
		"EclairGeoData",
		"EinsteinAssistantSettings",
		"EmailAdministrationSettings",
		"EmailFolder",
		"EmailIntegrationSettings",
		"EmailServicesFunction",
		"EmailTemplate",
		"EmailTemplateSettings",
		"EmbeddedServiceBranding",
		"EmbeddedServiceConfig",
		"EmbeddedServiceFlowConfig",
		"EmbeddedServiceLiveAgent",
		"EnhancedNotesSettings",
		"EntitlementProcess",
		"EntitlementSettings",
		"EntitlementTemplate",
		"EntityImplements",
		"EscalationRules",
		"EssentialsSettings",
		"EventSettings",
		"ExperienceBundle",
		"ExperienceBundleSettings",
		"ExternalDataSource",
		"ExternalServiceRegistration",
		"ExternalServicesSettings",
		"FeatureParameterBoolean",
		"FeatureParameterDate",
		"FeatureParameterInteger",
		"FieldServiceMobileExtension",
		"FieldServiceSettings",
		"FieldSet",
		"FieldSrcTrgtRelationship",
		"FilesConnectSettings",
		"FileUploadAndDownloadSecuritySettings",
		"FlexiPage",
		"Flow",
		"FlowCategory",
		"FlowDefinition",
		"FlowSettings",
		"ForecastingSettings",
		"FormulaSettings",
		"FunctionReference",
		"GatewayProviderPaymentMethodType",
		"GlobalValueSet",
		"GlobalValueSetTranslation",
		"GoogleAppsSettings",
		"Group",
		"HighVelocitySalesSettings",
		"HomePageComponent",
		"HomePageLayout",
		"Icon",
		"IdeasSettings",
		"IframeWhiteListUrlSettings",
		"InboundCertificate",
		"InboundNetworkConnection",
		"Index",
		"IndustriesManufacturingSettings",
		"IndustriesSettings",
		"InstalledPackage",
		"InventorySettings",
		"InvocableActionSettings",
		"IoTSettings",
		"IsvHammerSettings",
		"KeywordList",
		"KnowledgeSettings",
		"LanguageSettings",
		"Layout",
		"LeadConfigSettings",
		"LeadConvertSettings",
		"Letterhead",
		"LightningBolt",
		"LightningComponentBundle",
		"LightningExperienceSettings",
		"LightningExperienceTheme",
		"LightningMessageChannel",
		"LightningOnboardingConfig",
		"ListView",
		"LiveAgentSettings",
		"LiveChatAgentConfig",
		"LiveChatButton",
		"LiveChatDeployment",
		"LiveChatSensitiveDataRule",
		"LiveMessageSettings",
		"MacroSettings",
		"ManagedContentType",
		"ManagedTopics",
		"MapsAndLocationSettings",
		"MatchingRules",
		"MilestoneType",
		"MlDomain",
		"MobileApplicationDetail",
		"MobileSettings",
		"ModerationRule",
		"MutingPermissionSet",
		"MyDomainDiscoverableLogin",
		"MyDomainSettings",
		"NamedCredential",
		"NameSettings",
		"NavigationMenu",
		"Network",
		"NetworkBranding",
		"NotificationsSettings",
		"NotificationTypeConfig",
		"OauthCustomScope",
		"ObjectLinkingSettings",
		"ObjectSourceTargetMap",
		"OmniChannelSettings",
		"OpportunityInsightsSettings",
		"OpportunityScoreSettings",
		"OpportunitySettings",
		"OrderManagementSettings",
		"OrderSettings",
		"OrgSettings",
		"OutboundNetworkConnection",
		"PardotEinsteinSettings",
		"PardotSettings",
		"ParticipantRole",
		"PartyDataModelSettings",
		"PathAssistant",
		"PathAssistantSettings",
		"PaymentGatewayProvider",
		"PermissionSet",
		"PermissionSetGroup",
		"PicklistSettings",
		"PlatformCachePartition",
		"PlatformEventChannel",
		"PlatformEventChannelMember",
		"PortalsSettings",
		"PostTemplate",
		"PredictionBuilderSettings",
		"PresenceDeclineReason",
		"PresenceUserConfig",
		"PrivacySettings",
		"ProductSettings",
		"Profile",
		"ProfilePasswordPolicy",
		"ProfileSessionSetting",
		"Prompt",
		"Queue",
		"QueueRoutingConfig",
		"QuickAction",
		"QuickTextSettings",
		"QuoteSettings",
		"RecommendationBuilderSettings",
		"RecommendationStrategy",
		"RecordActionDeployment",
		"RecordPageSettings",
		"RecordType",
		"RedirectWhitelistUrl",
		"RemoteSiteSetting",
		"ReportFolder",
		"ReportType",
		"RestrictionRule",
		"RetailExecutionSettings",
		"Role",
		"SalesAgreementSettings",
		"SalesWorkQueueSettings",
		"SamlSsoConfig",
		"SchemaSettings",
		"SearchSettings",
		"SecuritySettings",
		"ServiceChannel",
		"ServiceCloudVoiceSettings",
		"ServicePresenceStatus",
		"ServiceSetupAssistantSettings",
		"SharingCriteriaRule",
		"SharingGuestRule",
		"SharingOwnerRule",
		"SharingReason",
		"SharingRules",
		"SharingSet",
		"SharingSettings",
		"SharingTerritoryRule",
		"SiteDotCom",
		"SiteSettings",
		"Skill",
		"SocialCustomerServiceSettings",
		"SocialProfileSettings",
		"SourceTrackingSettings",
		"StandardValue",
		"StandardValueSet",
		"StandardValueSetTranslation",
		"StaticResource",
		"SurveySettings",
		"SynonymDictionary",
		"SystemNotificationSettings",
		"Territory",
		"Territory2",
		"Territory2Model",
		"Territory2Rule",
		"Territory2Settings",
		"Territory2Type",
		"TimeSheetTemplate",
		"TrailheadSettings",
		"TransactionSecurityPolicy",
		"Translations",
		"TrialOrgSettings",
		"UIObjectRelationConfig",
		"UiPlugin",
		"UserAuthCertificate",
		"UserCriteria",
		"UserEngagementSettings",
		"UserInterfaceSettings",
		"UserManagementSettings",
		"UserProvisioningConfig",
		"ValidationRule",
		"WaveApplication",
		"WaveDashboard",
		"WaveDataflow",
		"WaveDataset",
		"WaveLens",
		"WaveRecipe",
		"WaveTemplateBundle",
		"WaveXmd",
		"WebLink",
		"WebStoreTemplate",
		"WebToXSettings",
		"WorkDotComSettings",
		"Workflow",
		"WorkflowAlert",
		"WorkflowFieldUpdate",
		"WorkflowFlowAction",
		"WorkflowKnowledgePublish",
		"WorkflowOutboundMessage",
		"WorkflowRule",
		"WorkflowSend",
		"WorkflowTask",
		"WorkSkillRouting",
	}
	// add support for only extracting certain objects
	if len(includeMetadataNames) > 0 {
		sort.Strings(includeMetadataNames)
		metadataNames = includeMetadataNames
	}

	for _, name := range metadataNames {
		if !isExcluded(name) {
			query = append(query, ForceMetadataQueryElement{Name: []string{name}, Members: []string{"*"}})
		}
	}

	if len(includeMetadataNames) == 0 {

		folders, err := force.GetAllFolders()
		if err != nil {
			err = fmt.Errorf("Could not get folders: %s", err.Error())
			ErrorAndExit(err.Error())
		}
		for foldersType, foldersName := range folders {
			if foldersType == "Email" {
				foldersType = "EmailTemplate"
			}
			members, err := force.GetMetadataInFolders(foldersType, foldersName)
			if err != nil {
				err = fmt.Errorf("Could not get metadata in folders: %s", err.Error())
				ErrorAndExit(err.Error())
			}

			if !isExcluded(string(foldersType)) {
				query = append(query, ForceMetadataQueryElement{Name: []string{string(foldersType)}, Members: members})
			}
		}
	}
	// fmt.Printf("Query: %s\n", query)

	if root == "" {
		root, err = config.GetSourceDir()
		if err != nil {
			fmt.Printf("Error obtaining root directory\n")
			ErrorAndExit(err.Error())
		}
	}
	files, problems, err := force.Metadata.Retrieve(query)
	if err != nil {
		fmt.Printf("Encountered and error with retrieve...\n")
		ErrorAndExit(err.Error())
	}
	if showWarnings {
		for _, problem := range problems {
			fmt.Fprintln(os.Stderr, problem)
		}
	}
	for name, data := range files {
		file := filepath.Join(root, name)
		dir := filepath.Dir(file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			ErrorAndExit(err.Error())
		}
		if err := ioutil.WriteFile(filepath.Join(root, name), data, 0644); err != nil {
			ErrorAndExit(err.Error())
		}
	}
	fmt.Printf("Exported to %s\n", root)
}

func isExcluded(name string) bool {
	index := sort.SearchStrings(excludeMetadataNames, name)

	return index < len(excludeMetadataNames) && excludeMetadataNames[index] == name
}

func isIncluded(name string) bool {
	index := sort.SearchStrings(includeMetadataNames, name)

	return index < len(includeMetadataNames) && includeMetadataNames[index] == name
}
