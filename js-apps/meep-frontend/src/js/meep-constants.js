/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export const HOST_PATH = location.origin;

// Sign In Status
export const STATUS_SIGNIN_NOT_SUPPORTED = 'SIGNED-NOT-SUPPORTED';
export const STATUS_SIGNED_IN = 'SIGNED-IN';
export const STATUS_SIGNING_IN = 'SIGNING-IN';
export const STATUS_SIGNED_OUT = 'SIGNED-OUT';
export const OAUTH_PROVIDER_GITHUB = 'github';
export const OAUTH_PROVIDER_GITLAB = 'gitlab';

// MEEP types
export const TYPE_CFG = 'CFG';
export const TYPE_EXEC = 'EXEC';

export const PAGE_CONFIGURE = 'PAGE_CONFIGURE';
export const PAGE_EXECUTE = 'PAGE_EXECUTE';
export const PAGE_MONITOR = 'PAGE_MONITOR';
export const PAGE_SETTINGS = 'PAGE_SETTINGS';
export const PAGE_HOME = 'PAGE_HOME';

// Page tab index
export const PAGE_HOME_INDEX = 0;
export const PAGE_CONFIGURE_INDEX = 1;
export const PAGE_EXECUTE_INDEX = 2;
export const PAGE_MONITOR_INDEX = 3;
export const PAGE_SETTINGS_INDEX = 4;

// URLs
export const MEEP_HELP_GUI_URL = 'https://interdigitalinc.github.io/AdvantEDGE/docs/usage/gui';
export const MEEP_ARCHITECTURE_URL = 'https://interdigitalinc.github.io/AdvantEDGE/docs/overview/overview-architecture';
export const MEEP_USAGE_URL = 'https://interdigitalinc.github.io/AdvantEDGE/docs/usage/usage-workflow';
export const MEEP_DOCS_URL = 'https://interdigitalinc.github.io/AdvantEDGE';
export const MEEP_GITHUB_URL = 'https://github.com/InterDigitalInc/AdvantEDGE';
export const MEEP_DISCUSSIONS_URL = 'https://github.com/InterDigitalInc/AdvantEDGE/discussions';
export const MEEP_LICENSE_URL = 'https://github.com/InterDigitalInc/AdvantEDGE/blob/master/LICENSE';
export const MEEP_CONTRIBUTING_URL = 'https://github.com/InterDigitalInc/AdvantEDGE/blob/master/CONTRIBUTING.md';
export const MEEP_ISSUES_URL = 'https://github.com/InterDigitalInc/AdvantEDGE/issues';

// MEEP IDs
export const MEEP_TAB_CFG = 'meep-tab-cfg';
export const MEEP_TAB_EXEC = 'meep-tab-exec';
export const MEEP_TAB_MON = 'meep-tab-mon';
export const MEEP_TAB_SET = 'meep-tab-set';
export const MEEP_TAB_HOME = 'meep-tab-home';
export const MEEP_LBL_SCENARIO_NAME = 'meep-lbl-scenario-name';
export const MEEP_BTN_CANCEL = 'meep-btn-cancel';
export const MEEP_BTN_APPLY = 'meep-btn-apply';

// Dialog IDs
export const MEEP_DLG_NEW_SANDBOX = 'meep-dlg-new-sandbox';
export const MEEP_DLG_NEW_SANDBOX_NAME = 'meep-dlg-new-sandbox-name';
export const MEEP_DLG_DELETE_SANDBOX = 'meep-dlg-delete-sandbox';
export const MEEP_DLG_NEW_SCENARIO = 'meep-dlg-new-scenario';
export const MEEP_DLG_NEW_SCENARIO_NAME = 'meep-dlg-new-scenario-name';
export const MEEP_DLG_SAVE_SCENARIO = 'meep-dlg-save-scenario';
export const MEEP_DLG_SAVE_REPLAY = 'meep-dlg-save-replay';
export const MEEP_DLG_SAVE_REPLAY_NAME = 'meep-dlg-save-replay-name';
export const MEEP_DLG_SAVE_REPLAY_DESCRIPTION = 'meep-dlg-save-replay-description';
export const MEEP_DLG_OPEN_SCENARIO = 'meep-dlg-open-scenario';
export const MEEP_DLG_OPEN_SCENARIO_SELECT = 'meep-dlg-open-scenario-select';
export const MEEP_DLG_DEL_SCENARIO = 'meep-dlg-del-scenario';
export const MEEP_DLG_INVALID_SCENARIO = 'meep-dlg-invalid-scenario';
export const MEEP_DLG_EXPORT_SCENARIO = 'meep-dlg-export-scenario';
export const MEEP_DLG_DEPLOY_SCENARIO = 'meep-dlg-deploy-scenario';
export const MEEP_DLG_DEPLOY_SCENARIO_SELECT = 'meep-dlg-deploy-scenario-select';
export const MEEP_DLG_TERMINATE_SCENARIO = 'meep-dlg-terminate-scenario';
export const MEEP_DLG_CONFIRM = 'meep-dlg-confirm';

// Dialog Types
// HOME
export const IDC_DIALOG_SIGN_IN = 'IDC_DIALOG_SIGN_IN';
export const IDC_DIALOG_SESSION_TERMINATED = 'IDC_DIALOG_SESSION_TERMINATED';
// CFG
export const IDC_DIALOG_OPEN_SCENARIO = 'IDC_DIALOG_OPEN_SCENARIO';
export const IDC_DIALOG_NEW_SCENARIO = 'IDC_DIALOG_NEW_SCENARIO';
export const IDC_DIALOG_SAVE_SCENARIO = 'IDC_DIALOG_SAVE_SCENARIO';
export const IDC_DIALOG_DELETE_SCENARIO = 'IDC_DIALOG_DELETE_SCENARIO';
export const IDC_DIALOG_EXPORT_SCENARIO = 'IDC_DIALOG_EXPORT_SCENARIO';
// EXEC
export const IDC_DIALOG_NEW_SANDBOX = 'IDC_DIALOG_NEW_SANDBOX';
export const IDC_DIALOG_DELETE_SANDBOX = 'IDC_DIALOG_DELETE_SANDBOX';
export const IDC_DIALOG_DEPLOY_SCENARIO = 'IDC_DIALOG_DEPLOY_SCENARIO';
export const IDC_DIALOG_TERMINATE_SCENARIO = 'IDC_DIALOG_TERMINATE_SCENARIO';
export const IDC_DIALOG_SAVE_REPLAY = 'IDC_DIALOG_SAVE_REPLAY';
// MONITORING
export const IDC_DIALOG_DELETE_DASHBOARD_LIST = 'IDC_DIALOG_DELETE_DASHBOARD_LIST';
// SETTINGS
export const IDC_DIALOG_CLEAR_UI_CACHE = 'IDC_DIALOG_CLEAR_UI_CACHE';

// Configuration page states
export const CFG_STATE_IDLE = 'IDLE';
export const CFG_STATE_NEW = 'NEW';
export const CFG_STATE_LOADED = 'LOADED';

// Configuration page views
export const CFG_VIEW_NETWORK = 'Network';
export const CFG_VIEW_MAP = 'Map';

// Configuration page IDs
export const CFG_VIS = 'cfg-vis';

export const CFG_VIEW_TYPE = 'cfg-view-type';

export const CFG_BTN_NEW_SCENARIO = 'cfg-btn-new-scenario';
export const CFG_BTN_OPEN_SCENARIO = 'cfg-btn-open-scenario';
export const CFG_BTN_SAVE_SCENARIO = 'cfg-btn-save-scenario';
export const CFG_BTN_DEL_SCENARIO = 'cfg-btn-del-scenario';
export const CFG_BTN_IMP_SCENARIO = 'cfg-btn-imp-scenario';
export const CFG_BTN_EXP_SCENARIO = 'cfg-btn-exp-scenario';
export const CFG_BTN_NEW_ELEM = 'cfg-btn-new-elem';
export const CFG_BTN_DEL_ELEM = 'cfg-btn-del-elem';
export const CFG_BTN_CLONE_ELEM = 'cfg-btn-clone-elem';
export const CFG_BTN_SAVE_ELEM = 'cfg-btn-save-elem';

export const CFG_ELEM_TYPE = 'cfg-elem-type';
export const CFG_ELEM_NAME = 'cfg-elem-name';
export const CFG_ELEM_PARENT = 'cfg-elem-parent';
export const CFG_ELEM_IMG = 'cfg-elem-img';
export const CFG_ELEM_GROUP = 'cfg-elem-group';
export const CFG_ELEM_ENV = 'cfg-elem-env';
export const CFG_ELEM_PORT = 'cfg-elem-port';
export const CFG_ELEM_EXT_PORT = 'cfg-elem-ext-port';
export const CFG_ELEM_PROT = 'cfg-elem-prot';
export const CFG_ELEM_GPU_COUNT = 'cfg-elem-gpu-count';
export const CFG_ELEM_GPU_TYPE = 'cfg-elem-gpu-type';
export const CFG_ELEM_CPU_MIN = 'cfg-elem-cpu-min';
export const CFG_ELEM_CPU_MAX = 'cfg-elem-cpu-max';
export const CFG_ELEM_MEMORY_MIN = 'cfg-elem-memory-min';
export const CFG_ELEM_MEMORY_MAX = 'cfg-elem-memory-max';
export const CFG_ELEM_PLACEMENT_ID = 'cfg-elem-placement-id';
export const CFG_ELEM_CMD = 'cfg-elem-cmd';
export const CFG_ELEM_ARGS = 'cfg-elem-args';
export const CFG_ELEM_EXTERNAL_CHECK = 'cfg-elem-external-check';
export const CFG_ELEM_MNC = 'cfg-elem-mnc';
export const CFG_ELEM_MCC = 'cfg-elem-mcc';
export const CFG_ELEM_MAC_ID = 'cfg-elem-mac-id';
export const CFG_ELEM_UE_MAC_ID = 'cfg-elem-ue-mac-id';
export const CFG_ELEM_DEFAULT_CELL_ID = 'cfg-elem-default-cell-id';
export const CFG_ELEM_CELL_ID = 'cfg-elem-cell-id';
export const CFG_ELEM_NR_CELL_ID = 'cfg-elem-nr-cell-id';
export const CFG_ELEM_GEO_LOCATION = 'cfg-elem-location';
export const CFG_ELEM_GEO_RADIUS = 'cfg-elem-radius';
export const CFG_ELEM_D2D_RADIUS = 'cfg-elem-d2d-radius';
export const CFG_ELEM_D2D_DISABLED = 'cfg-elem-d2d-disabled';
export const CFG_ELEM_GEO_PATH = 'cfg-elem-path';
export const CFG_ELEM_GEO_EOP_MODE = 'cfg-elem-eop-mode';
export const CFG_ELEM_GEO_VELOCITY = 'cfg-elem-velocity';
export const CFG_ELEM_CHART_CHECK = 'cfg-elem-chart-check';
export const CFG_ELEM_CHART_LOC = 'cfg-elem-chart-loc';
export const CFG_ELEM_CHART_GROUP = 'cfg-elem-chart-group';
export const CFG_ELEM_CHART_ALT_VAL = 'cfg-elem-chart-alt-val';
export const CFG_ELEM_CONNECTED = 'cfg-elem-connected';
export const CFG_ELEM_CONNECTIVITY_MODEL = 'cfg-elem-connectivity-model';
export const CFG_ELEM_DN_NAME = 'cfg-elem-dn-name';
export const CFG_ELEM_DN_LADN_CHECK = 'cfg-elem-dn-ladn-check';
export const CFG_ELEM_DN_ECSP = 'cfg-elem-dn-ecsp';
export const CFG_ELEM_WIRELESS = 'cfg-elem-wireless';
export const CFG_ELEM_WIRELESS_TYPE = 'cfg-elem-wireless-type';
export const CFG_ELEM_LATENCY = 'cfg-elem-latency';
export const CFG_ELEM_LATENCY_VAR = 'cfg-elem-latency-var';
export const CFG_ELEM_LATENCY_DIST = 'cfg-elem-latency-dist';
export const CFG_ELEM_PKT_LOSS = 'cfg-elem-pkt-loss';
export const CFG_ELEM_THROUGHPUT_DL = 'cfg-elem-throughput-dl';
export const CFG_ELEM_THROUGHPUT_UL = 'cfg-elem-throughput-ul';
export const CFG_ELEM_INGRESS_SVC_MAP = 'cfg-elem-ingress-svc-map';
export const CFG_ELEM_EGRESS_SVC_MAP = 'cfg-elem-egress-svc-map';
export const CFG_ELEM_META_DISPLAY_MAP_COLOR = 'cfg-elem-meta-display-map-color';


// Execution page states
export const EXEC_STATE_IDLE = 'IDLE';
export const EXEC_STATE_DEPLOYED = 'DEPLOYED';

// Execution page IDs
export const EXEC_SELECT_SANDBOX = 'exec-select-sandbox';

export const EXEC_BTN_NEW_SANDBOX = 'exec-btn-new-sandbox';
export const EXEC_BTN_DELETE_SANDBOX = 'exec-btn-delete-sandbox';
export const EXEC_BTN_DEPLOY = 'exec-btn-deploy';
export const EXEC_BTN_SAVE_SCENARIO = 'exec-btn-save-scenario';
export const EXEC_BTN_TERMINATE = 'exec-btn-terminate';
export const EXEC_BTN_DASHBOARD = 'exec-btn-dashboard';
export const EXEC_BTN_DASHBOARD_BTN_CLOSE = 'exec-btn-dashboard-btn-close';
export const EXEC_BTN_EVENT = 'exec-btn-event';
export const EXEC_BTN_EVENT_BTN_MANUAL_REPLAY = 'exec-btn-event-btn-manual-replay';
export const EXEC_BTN_EVENT_BTN_AUTOMATION = 'exec-btn-event-btn-automation';
export const EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOVEMENT = 'exec-btn-event-btn-automation-chkbox-movement';
export const EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_MOBILITY = 'exec-btn-event-btn-automation-chkbox-mobility';
export const EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_POAS_IN_RANGE = 'exec-btn-event-btn-automation-chkbox-poas-in-range';
export const EXEC_BTN_EVENT_BTN_AUTOMATION_CHKBOX_NETCHAR = 'exec-btn-event-btn-automation-chkbox-netchar';
export const EXEC_BTN_EVENT_BTN_AUTOMATION_BTN_CLOSE = 'exec-btn-event-btn-automation-btn-close';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY = 'exec-btn-event-btn-auto-replay';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_START = 'exec-btn-event-btn-auto-replay-btn-replay-start';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_REPLAY_STOP = 'exec-btn-event-btn-auto-replay-btn-replay-stop';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY_BTN_CLOSE = 'exec-btn-event-btn-auto-replay-btn-close';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY_CHKBOX_LOOP = 'exec-btn-event-btn-auto-replay-chkbox-loop';
export const EXEC_BTN_EVENT_BTN_AUTO_REPLAY_EVT_REPLAY_FILES = 'exec-btn-event-btn-auto-replay-evt-replay-files';
export const EXEC_BTN_EVENT_BTN_SAVE_REPLAY = 'exec-btn-event-btn-save-replay';
export const EXEC_BTN_EVENT_BTN_CLOSE = 'exec-btn-event-btn-close';

export const EXEC_VIEW_SELECT = 'exec-view-select';

export const EXEC_EVT_TYPE = 'exec-evt-type';
export const EXEC_EVT_MOB_TARGET = 'exec-evt-mob-target';
export const EXEC_EVT_MOB_DEST = 'exec-evt-mob-dest';
export const EXEC_EVT_NC_TYPE = 'exec-evt-nc-type';
export const EXEC_EVT_NC_NAME = 'exec-evt-nc-name';

export const EXEC_EVT_SU_ACTION = 'exec-evt-su-action';
export const EXEC_EVT_SU_REMOVE_ELEM_TYPE = 'exec-evt-su-remove-elem-type';
export const EXEC_EVT_SU_REMOVE_ELEM_NAME = 'exec-evt-su-remove-elem-name';

export const EXEC_EVT_PDU_SESSION_ACTION = 'exec-evt-pdu-session-action';
export const EXEC_EVT_PDU_SESSION_UE = 'exec-evt-pdu-session-ue';
export const EXEC_EVT_PDU_SESSION_ID = 'exec-evt-pdu-session-id';
export const EXEC_EVT_PDU_SESSION_DNN = 'exec-evt-pdu-session-dnn';

export const MEEP_EVENT_COUNT = 'meep-event-count';

// Trivia
export const NO_SCENARIO_NAME = 'NO_SCENARIO_NAME_12Q(*&HGHG___--9098';

export const DOMAIN_TYPE_STR = 'OPERATOR';
export const DOMAIN_CELL_TYPE_STR = 'OPERATOR-CELLULAR';
export const PUBLIC_DOMAIN_TYPE_STR = 'PUBLIC';
export const ZONE_TYPE_STR = 'ZONE';
export const COMMON_ZONE_TYPE_STR = 'COMMON';
export const NL_TYPE_STR = 'POA';
export const POA_TYPE_STR = 'POA';
export const POA_4G_TYPE_STR = 'POA-4G';
export const POA_5G_TYPE_STR = 'POA-5G';
export const POA_WIFI_TYPE_STR = 'POA-WIFI';
export const DEFAULT_NL_TYPE_STR = 'DEFAULT';
export const UE_TYPE_STR = 'UE';
export const FOG_TYPE_STR = 'FOG';
export const EDGE_TYPE_STR = 'EDGE';
export const CN_TYPE_STR = 'CN';
export const DC_TYPE_STR = 'DC';
export const MEC_SVC_TYPE_STR = 'MEC-SVC';
export const UE_APP_TYPE_STR = 'UE-APP';
export const EDGE_APP_TYPE_STR = 'EDGE-APP';
export const CLOUD_APP_TYPE_STR = 'CLOUD-APP';

export const ELEMENT_TYPE_SCENARIO = 'SCENARIO';
export const ELEMENT_TYPE_OPERATOR = 'OPERATOR';
export const ELEMENT_TYPE_OPERATOR_GENERIC = 'OPERATOR GENERIC';
export const ELEMENT_TYPE_OPERATOR_CELL = 'OPERATOR CELLULAR';
export const ELEMENT_TYPE_ZONE = 'ZONE';
export const ELEMENT_TYPE_POA = 'POA';
export const ELEMENT_TYPE_POA_GENERIC = 'POA GENERIC';
export const ELEMENT_TYPE_POA_4G = 'POA CELLULAR 4G';
export const ELEMENT_TYPE_POA_5G = 'POA CELLULAR 5G';
export const ELEMENT_TYPE_POA_WIFI = 'POA WIFI';
export const ELEMENT_TYPE_DC = 'DISTANT CLOUD';
export const ELEMENT_TYPE_CN = 'CORE NETWORK';
export const ELEMENT_TYPE_EDGE = 'EDGE';
export const ELEMENT_TYPE_FOG = 'FOG';
export const ELEMENT_TYPE_UE = 'TERMINAL';
export const ELEMENT_TYPE_MECSVC = 'MEC SERVICE';
export const ELEMENT_TYPE_UE_APP = 'TERMINAL APPLICATION';
export const ELEMENT_TYPE_EDGE_APP = 'EDGE APPLICATION';
export const ELEMENT_TYPE_CLOUD_APP = 'CLOUD APPLICATION';

// Default latencies per physical location type
export const DEFAULT_LATENCY_INTER_DOMAIN = 50;
export const DEFAULT_LATENCY_JITTER_INTER_DOMAIN = 10;
export const DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN = 'Normal';
export const DEFAULT_THROUGHPUT_DL_INTER_DOMAIN = 1000;
export const DEFAULT_THROUGHPUT_UL_INTER_DOMAIN = 1000;
export const DEFAULT_PACKET_LOSS_INTER_DOMAIN = 0;
export const DEFAULT_LATENCY_INTER_ZONE = 6;
export const DEFAULT_LATENCY_JITTER_INTER_ZONE = 2;
export const DEFAULT_THROUGHPUT_DL_INTER_ZONE = 1000;
export const DEFAULT_THROUGHPUT_UL_INTER_ZONE = 1000;
export const DEFAULT_PACKET_LOSS_INTER_ZONE = 0;
export const DEFAULT_LATENCY_INTRA_ZONE = 5;
export const DEFAULT_LATENCY_JITTER_INTRA_ZONE = 1;
export const DEFAULT_THROUGHPUT_DL_INTRA_ZONE = 1000;
export const DEFAULT_THROUGHPUT_UL_INTRA_ZONE = 1000;
export const DEFAULT_PACKET_LOSS_INTRA_ZONE = 0;
export const DEFAULT_LATENCY_TERMINAL_LINK = 1;
export const DEFAULT_LATENCY_JITTER_TERMINAL_LINK = 1;
export const DEFAULT_THROUGHPUT_DL_TERMINAL_LINK = 1000;
export const DEFAULT_THROUGHPUT_UL_TERMINAL_LINK = 1000;
export const DEFAULT_PACKET_LOSS_TERMINAL_LINK = 0;
export const DEFAULT_LATENCY_LINK = 0;
export const DEFAULT_LATENCY_JITTER_LINK = 0;
export const DEFAULT_THROUGHPUT_DL_LINK = 1000;
export const DEFAULT_THROUGHPUT_UL_LINK = 1000;
export const DEFAULT_PACKET_LOSS_LINK = 0;
export const DEFAULT_LATENCY_APP = 0;
export const DEFAULT_LATENCY_JITTER_APP = 0;
export const DEFAULT_THROUGHPUT_DL_APP = 1000;
export const DEFAULT_THROUGHPUT_UL_APP = 1000;
export const DEFAULT_PACKET_LOSS_APP = 0;
export const DEFAULT_LATENCY_DC = 0;

// Connection State & Types
export const OPT_CONNECTED = {label: 'Connected', value: true};
export const OPT_DISCONNECTED = {label: 'Disconnected', value: false};
export const OPT_WIRELESS = {label: 'Wireless', value: true};
export const OPT_WIRED = {label: 'Wired', value: false};

// Connectivity Models
export const CONNECTIVITY_MODEL_OPEN = 'OPEN';
export const CONNECTIVITY_MODEL_PDU = 'PDU';
export const DEFAULT_CONNECTIVITY_MODEL = CONNECTIVITY_MODEL_OPEN;

// D2d
export const DEFAULT_D2D_RADIUS = 100;
export const DEFAULT_D2D_DISABLED = false;

// GPU Types
export const GPU_TYPE_NVIDIA = 'NVIDIA';

// End-of-path modes
export const GEO_EOP_MODE_LOOP = 'LOOP';
export const GEO_EOP_MODE_REVERSE = 'REVERSE';

// Monitoring Page IDs
export const MON_DASHBOARD_SELECT = 'mon-dashboard-select';
export const MON_DASHBOARD_IFRAME = 'mon-dashboard-iframe';

// Settings Page IDs
export const SET_EXEC_REFRESH_CHECKBOX = 'set-exec-refresh-checkbox';
export const SET_EXEC_REFRESH_INT = 'set-exec-refresh-int';
export const SET_VIS_CFG_CHECKBOX = 'set-vis-cfg-checkbox';
export const SET_VIS_CFG_LABEL = 'VIS Configuration Mode';
export const SET_DASHBOARD_CFG_CHECKBOX = 'set-dashboard-cfg-checkbox';
export const SET_DASHBOARD_CFG_LABEL = 'Show Dashboard Config (Experimental)';
export const SET_RESET_SETTINGS_BUTTON = 'set-reset-settings-btn';

// Logical Scenario types
export const TYPE_SCENARIO = 0;
export const TYPE_DOMAIN = 1;
export const TYPE_ZONE = 2;
export const TYPE_NET_LOC = 3;
export const TYPE_PHY_LOC = 4;
export const TYPE_PROCESS = 5;

// NC Group Prefixes
export const PREFIX_INT_DOM = 'Inter-Domain';
export const PREFIX_INT_ZONE = 'Inter-Zone';
export const PREFIX_INTRA_ZONE = 'Intra-Zone';
export const PREFIX_TERM_LINK = 'Terminal Link';
export const PREFIX_LINK = 'Link';
export const PREFIX_APP = 'Application';

export const id = label => {
  return '#' + label;
};

export const VIEW_NAME_NONE = 'None';
export const MAP_VIEW = 'Map View';
export const NET_TOPOLOGY_VIEW = 'Network Topology';
export const SEQ_DIAGRAM_VIEW = 'Sequence Diagram';
export const DATAFLOW_DIAGRAM_VIEW = 'Data Flow Diagram';
export const NET_METRICS_PTP_VIEW = 'Network Metrics Point-to-Point';
export const NET_METRICS_AGG_VIEW = 'Network Metrics Aggregation';
export const WIRELESS_METRICS_PTP_VIEW = 'Wireless Metrics Point-to-Point';
export const WIRELESS_METRICS_AGG_VIEW = 'Wireless Metrics Aggregation';

export const DEST_DISCONNECTED = 'DISCONNECTED';

export const MOBILITY_EVENT = 'MOBILITY';
export const NETWORK_CHARACTERISTICS_EVENT = 'NETWORK-CHARACTERISTICS-UPDATE';
export const SCENARIO_UPDATE_EVENT = 'SCENARIO-UPDATE';
export const PDU_SESSION_EVENT = 'PDU-SESSION';

export const VIEW_1 = 'VIEW #1';
export const VIEW_2 = 'VIEW #2';

export const SCENARIO_UPDATE_ACTION_NONE = 'NONE';
export const SCENARIO_UPDATE_ACTION_ADD = 'ADD';
export const SCENARIO_UPDATE_ACTION_REMOVE = 'REMOVE';
export const SCENARIO_UPDATE_ACTION_MODIFY = 'MODIFY';

export const PDU_SESSION_ACTION_ADD = 'ADD';
export const PDU_SESSION_ACTION_REMOVE = 'REMOVE';

// Dashboard Config
export const DASH_SEQ_MAX_MSG_COUNT = 500;
export const DASH_DATAFLOW_MAX_MSG_COUNT = 10000;

// Default Dashboard list
export const DEFAULT_DASHBOARD_OPTIONS = [
  {
    label: 'None',
    value: ''
  },
  {
    label: 'Network Metrics Point-to-Point',
    value: HOST_PATH + '/grafana/d/1/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1s&theme=light<exec><vars>'
  },
  {
    label: 'Network Metrics Aggregation',
    value: HOST_PATH + '/grafana/d/2/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1s&theme=light<exec><vars>'
  },
  {
    label: 'Http REST API Logs Aggregation',
    value: HOST_PATH + '/grafana/d/3/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1s&theme=light'
  },
  {
    label: 'Http REST API Single Detailed Log',
    value: HOST_PATH + '/grafana/d/4/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1d&theme=light'
  },
  {
    label: 'Wireless Metrics Point-to-Point',
    value: HOST_PATH + '/grafana/d/5/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1s&theme=light<exec><vars>'
  },
  {
    label: 'Wireless Metrics Aggregation',
    value: HOST_PATH + '/grafana/d/6/metrics-dashboard?orgId=1&var-datasource=meep-influxdb&refresh=1s&theme=light<exec><vars>'
  },
  {
    label: 'Platform Metrics',
    value: HOST_PATH + '/grafana/d/platform-advantedge/platform-advantedge?orgId=1&refresh=15s&kiosk=tv&theme=light'
  },
  {
    label: 'Runtime Environment Metrics (Node)',
    value: HOST_PATH + '/grafana/d/runtime-environment-node/runtime-environment-node?orgId=1&refresh=15s&kiosk=tv&theme=light'
  }
];
