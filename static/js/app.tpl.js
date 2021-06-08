Vue.config.devtools = true;
const onError = function (error) {
  if (error.message === undefined) {
    if (error.response && error.response.status === 401) {
      this.onError = { message: "Unauthorized Access. Please check your token." }
    } else {
      this.onError = { message: "Something went wrong. Make sure the configuration is ok and your vermon Backend is up and running."}
    }
  } else {
    this.onError = { message: error.message }
  }
  console.log(this.onError.message)
}

function lastRun() {
  return moment().format('ddd, YYYY-MM-DD HH:mm:ss')
}

function theTime(tz){
  return moment().tz(tz).format('HH:mm:ss');
}

function getParam(parameter, defParam){
  theParam = getParameterByName(parameter);
  if (theParam == null)
    theParam=defParam
    return theParam
}

// Used by vue
const app = new Vue({
  el: '#app',
  data: {
    project: null,
    projects: {},
    servers: [],
    serversMap: {},
    token: null,
    vermonBack: null,
    loading: false,
    invalidConfig: false,
    lastRun: lastRun(),
    timezone1 : "Europe/Brussels",
    meTime : theTime("Europe/Brussels"),
	  count: 0,
    onError: null
  },
  created: function() {
    this.loadConfig()
    this.setupDefaults()
    console.log(axios.defaults.baseURL)
    this.fetchProject()
    var self = this
	  var count=self.refresh
    setInterval(function() {
    	if (count-- == 0){
        self.updateBuilds()
    	  count=self.refresh;;
    	}
    	self.count=count;
      self.meTime = theTime(self.timezone1);
    }, 1000)
  },
  methods: {
    loadConfig: function() {
      const self = this
      self.vermonBack = getParameterByName("vermonBack")
      if (self.vermonBack == null)
         self.vermonBack = "localhost:3000"
      self.proto = "http"
	    self.refresh = getParameterByName("refresh")
      self.timezone1 = getParameterByName("timezone1")
      if (self.timezone1 == null)
         self.timezone1 ="Europe/Brussels"
	    if (self.refresh == null)
		    self.refresh = 120;
	    self.count=self.refresh
      const project = getParameterByName("projects")
      self.project = project
    },
    /*
    setupDefaults: function() {
      if (this.token !== "use_cookie") {
        port="###PORT###"
        axios.defaults.baseURL = "http://localhost:"+port+ "/api/v1";
        axios.defaults.baseURL = this.proto + "://" + this.vermonBack + "/api/v1"
      } else {
        // Running on the vermonBack-Server...
        axios.defaults.baseURL = "/api/v1"
        this.vermonBack = location.hostname
      }
    },
    */
     setupDefaults: function () {
      port="###PORT###"
      axios.defaults.baseURL = "http://localhost:"+port+"/api/v1";
    },
    fetchProject: function() {
      const self = this
        self.loading = true
        request='project/'+self.project+'/getProject'
        console.log(request)
        axios.get(request)
          .then(function (response) {
            self.loading = false
            const project = response.data
            self.fetchVersions(project)
          })
          .catch(onError.bind(self))
    },
    updateBuilds: function() {
      const self = this
      self.onError = null
      Object.values(self.project).forEach(function(p) { 
          //self.fetchVersions(p)
          location.reload(); // big fat and dirty workaround
        })
      self.lastRun = lastRun()
    },
    fetchVersions: function(p) {
      
      const self = this
      axios.get('/project/' + p[0].Projectid + '/getVersions')
        .then(function(versions) {
          Object.values(versions.data).forEach(function(p) {
            self.updateBuildInfo(p)
          })
        })
        .catch(onError.bind(self))
      },
    updateBuildInfo: function(p) {
      const self = this
      key = p.AppId + "" + p.ServerId
      const apps = {
          key: key,
          AppDesc: p.AppDesc.String,
          AppId: p.AppId,
          AppName: p.AppName,
          AppVer: p.AppVer,
          IsSubApp: p.IsSubApp,
          PrjID: p.PrjID,
          ProjectDesc: p.ProjectDesc.String,
          ProjectName: p.ProjectName,
          ProjectTeam: p.ProjectTeam,
          Projectid: p.Projectid,
          ServerDesc: p.ServerDesc.String,
          ServerHName: p.ServerHName,
          ServerId: p.ServerId,
          ServerLName: p.ServerLName,
          parentApp: p.parentApp,
          releaseDate: p.releaseDate,
          PipelineId: p.PipelineId,
          PipelineUrl: p.PipelineUrl.String,
          ServerTier: p.ServerTier.String,
          ReleaseNote: p.ReleaseNote.String
        }
      self.servers.push(apps)
      self.serversMap[key] = apps
      //self.servers.sort(function(a, b) { return a.project.localeCompare(b.project) })
    }
  }
})