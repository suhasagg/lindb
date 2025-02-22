/*
Licensed to LinDB under one or more contributor
license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright
ownership. LinDB licenses this file to you under
the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
const local = {
  SiderMenu: {
    Overview: "概览",
    Configuration: "配置",
    Search: "查询",
    Explore: "浏览",
    Monitoring: "自监控",
    Dashboard: "监控看板",
    Metadata: "元数据",
    Replication: "副本复制",
    Request: "请求",
    "Log View": "日志",
    Database: "数据库",
    "Logic Database": "逻辑库",
    Storage: "存储集群",
    Broker: "计算集群",
    "Multiple IDCs": "多机房",
  },
  LayoutHeader: {
    language: "中文",
  },
  Overview: {
    brokerLiveNodes: "计算节点(Broker live nodes)",
  },
  NodeView: {
    hostIp: "主机 IP",
    hostName: "主机名",
    grpcPort: "GRPC 端口",
    httpPort: "HTTP 端口",
    title: "主机信息",
    uptime: "在线时间",
    version: "版本",
    cpu: "CPU",
    memory: "内存",
    nodeId: "节点 ID",
  },
  StorageView: {
    name: "集群名(Namespace)",
    nodeStatus: "节点状态",
    numOfDatabase: "数据库个数",
    replicationStatus: "副本状态",
    diskCapacityUsage: "磁盘使用",
    storageClusterList: "存储集群",
    totalOfReplication: "总副本",
    underReplicated: "在线",
    unavailableReplica: "离线",
    aliveNodes: "在线",
    deadNodes: "离线",
    liveNodes: "节点列表(Live Node)",
    databaseList: "数据库列表",
  },
  DatabaseView: {
    name: "名称",
    numOfShards: "分片数(Num. Of Shard)",
    replicaFactor: "副本因子(Replica Factor)",
  },
  CapacityView: {
    total: "总量",
    used: "已用",
    free: "空闲",
  },
  MasterView: {
    master: "调度节点(Master)",
    electTime: "选举时间",
  },
  SearchView: {
    database: "数据库名",
    databaseRequired: "请选择数据库名",
    search: "查询",
  },
  DataExploreView: {
    database: "数据库",
    namespace: "命名空间",
    metric: "指标名",
    showLinQL: "显示 LinQL",
    metricRequired: "请选择指标名",
    field: "字段",
    filterBy: "过滤",
    groupBy: "Group By",
  },
  ReplicationView: {
    database: "数据库",
    replicationStatus: "复制状态",
    shard: "分片",
    append: "写入",
    consume: "消费",
    ack: "持久化",
    lag: "积压",
    peer: "通道",
    type: "类型",
    memoryDatabaseStatus: "内存数据库状态",
    node: "节点",
    state: "状态",
    uptime: "创建时间",
    memSize: "内存大小",
    numOfMetrics: "指标数",
    numOfSeries: "Series 数",
    noMemoryDatabase: "暂无活跃的内存数据库",
    replica: "复制",
    memoryDatabase: "内存数据库",
  },
  LogView: {
    role: "角色",
    storage: "存储集群",
    node: "节点",
    file: "文件",
    size: "大小",
  },
  RequestView: {
    timestamp: "执行时间",
    duration: "执行耗时",
    linQL: "LinQL",
    database: "数据库",
    broker: "执行计算节点",
    runLinQL: "执行 LinQL",
  },
  MetadataExploreView: {
    compare: "对比",
    comparing: "正在进行对比",
    compareResult1: "共找到",
    compareResult2: "个节点，其中存在",
    compareResult3: "个不相同的节点",
    compareTooltip: "与内存中状态机中的数据进行对比",
    compareResultTitle: "状态对比结果",
    compareResultDesc: "持久化存储中的状态与内存状态品中的状态进行对比",
    filterNode: "节点过滤",
  },
  MetadataDatabaseView: {
    name: "数据库名",
    nameRequired: "请输入数据库名",
    storage: "存储集群",
    storageRequired: "请选择存储集群",
    description: "描述",
    deleteConfirm1: "您确认要删除数据库 [ ",
    deleteConfirm2: " ] ?",
    numOfShards: "分片数(Num. Of Shard)",
    numOfShardsRequired: "请输入分片数(Num. Of Shard)",
    replicaFactor: "副本因子(Replica Factor)",
    replicaFactorRequired: "请输入副本因子(Replica Factor)",
    engineOptions: "存储引擎配置",
    autoCreateNS: "自动创建命名空间",
    intervals: "存储间隔",
    interval: "存储间隔(秒)",
    retention: "存储时长(天)",
    writeableTimeRange: "可接受写入时间",
    behead: "往前",
    ahead: "往后",
    example: "例如: [ now()-1h ~ now()+1h ]",
  },
  MetadataLogicDatabaseView: {
    name: "数据库名",
    nameRequired: "请输入数据库名",
    router: "路由",
    tagKey: "路由键(Tag Key)",
    tagValues: "路由值(Tag Values)",
    brokers: "计算集群",
    deleteConfirm1: "您确认要删除数据库 [ ",
    deleteConfirm2: " ] ?",
  },
  MetadataClusterView: {
    register: "注册",
    name: "集群名(Namespace)",
    status: "状态",
    configuration: "配置",
    Ready: "在线",
    Initialize: "初始化",
    endpoints: "地址(ETCD)",
    username: "用户名(ETCD)",
    password: "密码(ETCD)",
    dialTimeout: "连接超时时间",
    timeout: "客户端处理超时时间",
    timeoutTooltip: "客户端处理请求的超时时间",
    leaseTTL: "租约(Time To Live)",
  },
  MetadataStorageView: {
    recoverConfirmMessage: "是否确认从存储集群的本地存储中恢复数据库配置？",
    recoverErrorTitle: "恢复数据库配置失败",
    recoverSuccessTitle: "恢复数据库配置成功",
  },
  LinSelectView: {
    placeholder: "请选择",
  },
  TimePicker: {
    from: "开始时间",
    to: "结束时间",
    searchQuickRange: "过滤时间区间",
    applyTimeRange: "确定",
    absoluteTimeRange: "固定时间区间",
    last15Min: "最近15分钟",
    last30Min: "最近30分钟",
    last1Hour: "最近1小时",
    last3Hour: "最近3小时",
    last6Hour: "最近6小时",
    last12Hour: "最近12小时",
    last1Day: "最近1天",
    last2Day: "最近2天",
    last3Day: "最近3天",
    last7Day: "最近7天",
    last15Day: "最近15天",
    last30Day: "最近30天",
    off: "关闭",
    "10s": "10秒",
    "30s": "30秒",
    "1m": "1分钟",
    "5m": "5分钟",
  },
  Common: {
    noData: "暂无数据",
    ok: "确定",
    cancel: "取消",
    submit: "提交",
    actions: "操作",
    create: "新建",
    pleaseConfirm: "请确认",
    unknownInternalError: "未知的内部错误",
    loading: "加载中",
  },
};

export default local;
