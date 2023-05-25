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
  -p, -package   # Include managed package metadata
  -n, -newline   # Add newline for output

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
	addNewline		   	   bool
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
	cmdExport.Flag.BoolVar(&addNewline, "n", false, "add newline")
	cmdExport.Flag.BoolVar(&addNewline, "newline", false, "add newline")
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
	standardValueSet := "StandardValueSet"

	sort.Strings(excludeMetadataNames)
	sort.Strings(includeMetadataNames)

	exportAll := true
	if len(includeMetadataNames) > 0 {
		exportAll = false
	}

	if (!isExcluded(customObject) && exportAll) || isIncluded(customObject) {
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
		stdObjects = append(stdObjects, "PersonAccount")

		query = append(query, ForceMetadataQueryElement{Name: []string{customObject}, Members: stdObjects})
	}

	standardValueSetNames := []string{
		"AccountContactMultiRoles",
		"AccountContactRole",
		"AccountOwnership",
		"AccountRating",
		"AccountType",
		"ACInitSumEmployeeType4",
		"ACInitSumInitiativeType4",
		"ACISumRecipientCategory4",
		"ACorruptionInitSumCountry4",
		"ACorruptionInitSumRegion4",
		"AssetActionCategory",
		"AssetRelationshipType",
		"AssetStatus",
		"AssociatedLocationType",
		"CampaignMemberStatus",
		"CampaignStatus",
		"CampaignType",
		"CardType",
		"CareItemStatus1",
		"CaseContactRole",
		"CaseOrigin",
		"CasePriority",
		"CaseReason",
		"CaseStatus",
		"CaseType",
		"CCPAdditionalBenefits4",
		"CCProjectMitigationType4",
		"CCPStandardsAgencyName4",
		"CCreditProjectProjectType4",
		"ChangeRequestRelatedItemImpactLevel",
		"ChangeRequestBusinessReason",
		"ChangeRequestCategory",
		"ChangeRequestImpact",
		"ChangeRequestPriority",
		"ChangeRequestRiskLevel",
		"ChangeRequestStatus",
		"CompanyRelationshipType4",
		"ContactPointAddressType",
		"ContactPointUsageType",
		"ContactRequestReason",
		"ContactRequestStatus",
		"ContactRole",
		"ContractContactRole",
		"ContractStatus",
		"ConsequenceOfFailure",
		"DigitalAssetStatus",
		"DEInclSumDiversityType4",
		"DEInclSumEmployeeType4",
		"DEInclSumEmploymentType4",
		"DEInclSumGender4",
		"DEISumDiversityCategory4",
		"DivrsEquityInclSumLocation4",
		"DivrsEquityInclSumRace4",
		"EntitlementType",
		"EBSEmployeeBenefitType4",
		"EBSPercentageCalcType4",
		"EBSummaryBenefitUsage4",
		"EBSummaryEmploymentType4",
		"EDemographicSumAgeGroup4",
		"EDemographicSumGender4",
		"EDemographicSumRegion4",
		"EDemographicSumReportType4",
		"EDemographicSumWorkType4",
		"EDevelopmentSumGender4",
		"EDSumEmployeeType4",
		"EDSumEmploymentType4",
		"EDSumProgramCategory4",
		"EPSumMarket4",
		"EPSumPerformanceCategory4",
		"EPSumPerformanceType4",
		"EPSumRegion4",
		"ERCompanyBusinessRegion4",
		"ERCompanySector4",
		"EReductionTargetTargetType4",
		"ERTargetOtherTargetKpi4",
		"ERTTargetSettingMethod4",
		"EventSubject",
		"EventType",
		"FinanceEventAction",
		"FinanceEventType",
		"FiscalYearPeriodName",
		"FiscalYearPeriodPrefix",
		"FiscalYearQuarterName",
		"FiscalYearQuarterPrefix",
		"ForecastingItemCategory2",
		"FreightHaulingMode4",
		"FtprntAuditApprovalStatus4",
		"FulfillmentStatus",
		"FulfillmentType",
		"GovtFinancialAsstSumType4",
		"IdeaCategory3",
		"IdeaMultiCategory",
		"IdeaStatus",
		"IdeaThemeStatus",
		"IncidentCategory",
		"IncidentImpact",
		"IncidentPriority",
		"IncidentRelatedItemImpactLevel",
		"IncidentRelatedItemImpactType",
		"IncidentReportedMethod",
		"IncidentStatus",
		"IncidentSubCategory",
		"IncidentType",
		"IncidentUrgency",
		"Industry",
		"LeadSource",
		"LeadStatus",
		"LocationType",
		"MilitaryService",
		"OIncidentSummaryHazardType4",
		"OISCorrectiveActionType4",
		"OISummaryIncidentSubtype4",
		"OISummaryIncidentType4",
		"OISummaryPenaltyType4",
		"OpportunityCompetitor",
		"OpportunityStage",
		"OpportunityType",
		"OrderItemSummaryChgRsn",
		"OrderStatus",
		"OrderSummaryRoutingSchdRsn",
		"OrderSummaryStatus",
		"OrderType",
		"PartnerRole",
		"PEFEFctrDataSourceType4",
		"ProblemCategory",
		"ProblemImpact",
		"ProblemPriority",
		"ProblemRelatedItemImpactLevel",
		"ProblemRelatedItemImpactType",
		"ProblemStatus",
		"ProblemSubCategory",
		"ProblemUrgency",
		"ProcessExceptionCategory",
		"ProcessExceptionPriority",
		"ProcessExceptionSeverity",
		"ProcessExceptionStatus",
		"Product2Family",
		"ProdRequestLineItemStatus",
		"ProductRequestStatus",
		"QuantityUnitOfMeasure",
		"QuestionOrigin3",
		"QuickTextCategory",
		"QuickTextChannel",
		"QuoteStatus",
		"RoleInTerritory2",
		"ResourceAbsenceType",
		"ReturnOrderLineItemProcessPlan",
		"ReturnOrderLineItemReasonForRejection",
		"ReturnOrderLineItemReasonForReturn",
		"ReturnOrderLineItemRepaymentMethod",
		"ReturnOrderShipmentType",
		"ReturnOrderStatus",
		"SalesTeamRole",
		"Salutation",
		"ScorecardMetricCategory",
		"ScienceBasedTargetStatus4",
		"SContributionSumCategory4",
		"Scope3CrbnFtprntStage4",
		"ServiceAppointmentStatus",
		"ServiceContractApprovalStatus",
		"ServTerrMemRoleType",
		"ShiftStatus",
		"SocialContributionSumType4",
		"SocialPostClassification",
		"SocialPostEngagementLevel",
		"SocialPostReviewedStatus",
		"SolutionStatus",
		"SourceBusinessRegion4",
		"StatusReason",
		"StnryAssetCrbnFtprntStage4",
		"StnryAstCrbnFtAllocStatus4",
		"StnryAstCrbnFtDataGapSts4",
		"StnryAstEvSrcStnryAstTyp4",
		"StnryAssetWaterFtprntStage4",
		"SupplierClassification4",
		"SupplierEmssnRdctnCmtTypev",
		"SupplierReportingScope4",
		"SupplierTier4",
		"SustainabilityScorecardStatus4",
		"TaskPriority",
		"TaskStatus",
		"TaskSubject",
		"TaskType",
		"UnitOfMeasure",
		"VehicleAstCrbnFtprntStage4",
		"VehicleType4",
		"WasteFootprintStage4",
		"WasteDisposalType4",
		"WasteType4",
		"WorkOrderLineItemPriority",
		"WorkOrderLineItemStatus",
		"WorkOrderPriority",
		"WorkOrderStatus",
		"WorkStepStatus",
		"WorkTypeDefApptType",
		"WorkTypeGroupAddInfo"
	}

	if (!isExcluded(standardValueSet) && exportAll) || isIncluded(standardValueSet) {
		query = append(query, ForceMetadataQueryElement{Name: []string{standardValueSet}, Members: standardValueSetNames})
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
		//"EmailTemplate",
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
		"PlatformEventSubscriberConfig",
		"Portal",
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
		"RestrictionRule",
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
		// "StandardValueSet",
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
		// "UIObjectRelationConfig",
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
		metadataNames = includeMetadataNames
	}

	for _, name := range metadataNames {
		if !isExcluded(name) {
			query = append(query, ForceMetadataQueryElement{Name: []string{name}, Members: []string{"*"}})
		}
	}

	if (exportAll || isIncluded("EmailTemplate"))  {

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
	newline := []byte("")
	if addNewline {
		newline = []byte("\n")		
	}
	for name, data := range files {
		file := filepath.Join(root, name)
		dir := filepath.Dir(file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			ErrorAndExit(err.Error())
		}
		if err := ioutil.WriteFile(filepath.Join(root, name), append(data, newline...), 0644); err != nil {
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
