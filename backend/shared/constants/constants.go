package constants

const (
	EventTopics = "EVENT_TOPICS"

	UacdEvent       = "uacd"
	UacdServiceName = "uacd"

	PromptEvent         = "prompt"
	PromptCanceledEvent = "prompt::cancel"
	StdoutOnlyLog       = "STD_OUT_ONLY"

	SipService                  = "sip_service"
	TelephoneService            = "telephone"
	ChatService                 = "chat"
	VideochatService            = "video_chat"
	VideoChatChannelType        = "video-chat"
	AssecoChatService           = "asseco-chat"
	EmailService                = "email"
	SmsService                  = "sms"
	ReportingService            = "reporting"
	TicketingService            = "ticketing"
	IamService                  = "iam"
	SipServiceCallModel         = "connection"
	SipServiceExtensionModel    = "extension"
	SipServiceCallChannelOpened = "opened"
	SipServiceCallInitiated     = "initiated"
	SipServiceCallRinging       = "ringing"
	SipServiceCallAnswered      = "answered"
	SipServiceCallHungUp        = "hangup"
	SipServiceCallFailed        = "failed"
	//conferenc actions
	JoinedConference = "conference::joined"
	LeftConference   = "conference::left"

	SipServiceExtensionRegistered   = "registered"
	SipServiceExtensionUnRegistered = "un-registered"

	ActivityCreatedAction = "activity::created"
	ActivityUpdatedAction = "activity::updated"
	ActivityDeletedAction = "activity::deleted"

	ServiceName           = "stem24_backend"
	AgentModel            = "agent"
	ActivityModel         = "activity"
	SkillGroupModel       = "skill_group"
	SkillGroupMemberModel = "skill_group_member"
	RequeueAction         = "queue::requeue"
	RequeuedAction        = "queue::requeued"
	QueuedAction          = "queue::queued"
	MatchAction           = "queue::match"
	MatchedAction         = "queue::matched"
	AnswerAction          = "queue::answer"
	AnsweredAction        = "queue::answered"
	DeletedQueueAction    = "queue::deleted"
	DequeuedQueueAction   = "queue::dequeued"
	ChannelModel          = "channel"
	CreatedAction         = "created"
	DeletedAction         = "deleted"
	UpdatedAction         = "updated"
	StatusChangeAction    = "status_change"
	RealmModel            = "realm"

	AgentJoinedSkillGroupAction = "skill_group::joined"
	AgentLeftSkillGroupAction   = "skill_group::left"

	NotificationTopic = "agent_management_notifications"
	DwhTopic          = "dwh"
	TelephoneChannel  = "TELEPHONE_CHANNEL"
	EmailChannel      = "EMAIl_CHANNEL"
	SmsChannel        = "SMS_CHANNEL"
	VideoChatChannel  = "VIDEO_CHAT_CHANNEL"
	AssecoChatChannel = "ASSECO_CHAT_CHANNEL"

	LogLevel = "LOG_LEVEL"

	GraceLoginPeriod = "GRACE_LOGIN_PERIOD"

	AgentNamespace = "/"

	XCorrelationID   = "X-Correlation-ID"
	TenantIdentifier = "X-TENANT-ID"
	DefaultTenant    = "default"
	Method           = "method"

	SingleTenantMode = "SINGLE_TENANT_MODE"

	SipEvent            = "sip"
	ActivityUpdateEvent = "ACTIVITY::update"

	ConsultationTypeConferenceType = "consultation"

	PersistPonderSnapshots          = "PERSIST_PONDER_SNAPSHOTS"
	IdlePonderTimeSpan              = "IDLE_PONDER_TIME_SPAN"
	IdlePonderTimeSpanDefault       = 60
	ActivitiesPonderTimeSpan        = "ACTIVITIES_PONDER_TIME_SPAN"
	ActivitiesPonderTimeSpanDefault = 60

	QueueResolverSleepTime        = "QUEUE_RESOLVE_SLEEP_TIME"
	QueueResolverSleepTimeDefault = 4000

	DirectoryUri         = "DIRECTORY_URI"
	DefaultDirectoryUri  = "http://live-directory:8080/v1/directory/default"
	TelephoneUri         = "TELEPHONE_URI"
	DefaultTelephoneUri  = "http://telephone:8080"
	VideoChatUri         = "VIDEO_CHAT_URI"
	DefaultVideoChatUri  = "http://video-chat:8080"
	AssecoChatUri        = "ASSECO_CHAT_URI"
	DefaultAssecoChatUri = "http://social-networks:8080"

	ActivityStatusHistoryEnabled = "ACTIVITY_STATUS_HISTORY_ENABLED"
	AgentStatusHistoryEnabled    = "AGENT_HISTORY_ENABLED"
	MovePreferredAgentToQueue    = "MOVE_PREFERRED_AGENT_TO_QUEUE"

	MaxUtilization        = "MAX_UTILIZATION"
	DefaultMaxUtilization = 100

	ShouldRunMigrations        = "RUN_MIGRATIONS"
	ShouldRunMigrationsDefault = true

	SkillGroupCronSchedule         = "SKILL_GROUP_CRON_SCHEDULE"
	WorkersPerSystemChannel        = "WORKERS_PER_SYSTEM_CHANNEL"
	DefaultWorkersPerSystemChannel = 3

	DefaultPerPage = 15
	DefaultPage    = 0

	QueuedItemSourceApi    = "api"
	QueuedItemSourceBroker = "broker"

	ActivityUuid   = "activity_uuid"
	SessionUuid    = "session_uuid"
	ConferenceUuid = "conference_uuid"
	Model          = "model"
	Action         = "action"
	Agent          = "agent"
	Contact        = "contact"
	QueuedItemUuid = "queued_item_uuid"
	Service        = "service"
	AgentId        = "agent_id"
	Number         = "number"
	ChannelId      = "channel_id"
	Extension      = "extension"
	ChannelType    = "channel_type"

	AppName = "APP_NAME"

	ClientId     = "CLIENT_ID"
	ClientSecret = "CLIENT_SECRET"

	IamUri             = "IAM_URI"
	IamRealm           = "IAM_REALM"
	PubKeySaveLocation = "PUB_KEY_LOCATION"
	Host               = "HOST"
	Port               = "PORT"

	BrokerType   = "BROKER_TYPE"
	BrokerHost   = "BROKER_HOST"
	BrokerPort   = "BROKER_PORT"
	BrokerUser   = "BROKER_USER"
	BrokerPass   = "BROKER_PASS"
	BrokerTopics = "BROKER_TOPICS"

	DatabaseDriver = "DB_DRIVER"
	DatabaseHost   = "DB_HOST"
	DatabasePort   = "DB_PORT"
	DatabaseUser   = "DB_USER"
	DatabaseName   = "DB_NAME"
	DatabasePass   = "DB_PASS"
	DebugDatabase  = "DB_DEBUG"
	DatabaseSeed   = "DB_SEED"

	ContentType     = "Content-Type"
	JsonContentType = "application/json"

	UserDisplayName = "display_name"
	UserId          = "user_id"
	AnyTenant       = "default"

	FileSystemDriver = "FILESYSTEM_DRIVER"

	ContentUrl        = "CONTENT_API_URL"
	ContentRoot       = "CONTENT_ROOT"
	ContentRepository = "CONTENT_REPOSITORY"

	LogDebugLevel = "DEBUG"
	LogInfoLevel  = "INFO"
	LogWarnLevel  = "WARN"
	LogErrorLevel = "ERROR"

	ServiceToken = "SERVICE_TOKEN"
	RefreshToken = "REFRESH_TOKEN"

	RecoveryTimer = "RECOVERY_TIMER"
	RecoveryCount = "RECOVERY_COUNT"

	MaxIdleConns    = "MAX_IDLE_CONNS"
	MaxOpenConns    = "MAX_OPEN_CONNS"
	ConnMaxLifetime = "CONN_MAX_LIFETIME"

	GraylogPort     = "GRAYLOG_PORT"
	GraylogHostname = "GRAYLOG_HOSTNAME"

	QueueLog        = "QUEUE_LOG"
	RequestLog      = "REQUEST_LOG"
	QueryLog        = "QUERY_LOG"
	ResponseBodyLog = "RESPONSE_BODY_LOG"
	LogDrivers      = "LOG_DRIVERS"
	LogGraylog      = "graylog"
	LogStdOut       = "stdout"
	LogFile         = "file"

	BrokerTopicDurable              = "1"
	BrokerTopicReadFromServiceStart = "1"
	BrokerTopicsSeparator           = ","
	BrokerTopicComponentSeparator   = "::"
	BrokerMessageCreationTime       = "CreationTime"

	RecordNotFound      = "record not found"
	Forbidden           = "forbidden"
	Conflict            = "conflict"
	BadRequest          = "bad request"
	UnprocessableEntity = "unprocessable entity"
	Internal            = "internal"

	RecordNotFoundCode      = "404"
	ForbiddenCode           = "403"
	ConflictCode            = "409"
	BadRequestCode          = "400"
	UnprocessableEntityCode = "422"
	InternalCode            = "500"
)
