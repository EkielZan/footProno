Vue.config.devtools = true;
const onError = function (error) {
  if (error.message === undefined) {
    if (error.response && error.response.status === 401) {
      this.onError = {
        message: "Unauthorized Access. Please check your token.",
      };
    } else {
      this.onError = {
        message:
          "Something went wrong. Make sure the configuration is ok and your vermon Backend is up and running.",
      };
    }
  } else {
    this.onError = { message: error.message };
  }
  console.log(this.onError.message);
};

function lastRun() {
  return moment().format("ddd, YYYY-MM-DD HH:mm:ss");
}

function theTime(tz) {
  return moment().tz(tz).format("HH:mm:ss");
}

function getParam(parameter, defParam) {
  theParam = getParameterByName(parameter);
  if (theParam == null) theParam = defParam;
  return theParam;
}

// Used by vue
const app = new Vue({
  el: "#app",
  data: {
    players: {},
    stats: {},
    backend: null,
    loading: false,
    invalidConfig: false,
    lastRun: lastRun(),
    timezone1: "Europe/Brussels",
    meTime: theTime("Europe/Brussels"),
    count: 0,
    onError: null,
    byPlayer: false,
  },
  created: function () {
    this.loadConfig();
    this.setupDefaults();
    console.log(axios.defaults.baseURL);
    var self = this;
    if (self.byPlayer) {
      this.fetchPlayer();
    } else {
      this.fetchPlayers();
    }
    this.fetchStats();
    setInterval(function () {
      if (count-- == 0) {
        if (self.byPlayer) {
          this.fetchPlayer();
        } else {
          this.fetchPlayers();
        }
        this.fetchStats();
        count = self.refresh;
      }
      self.count = count;
      self.meTime = theTime(self.timezone1);
    }, 1000);
  },
  methods: {
    loadConfig: function () {
      const self = this;
      self.backend = "localhost:8090";
      self.proto = "https";
      self.timezone1 = "Europe/Brussels";
      self.refresh = getParameterByName("refresh")
      if (self.refresh == null)
        self.refresh = 120;
      self.playerId = getParameterByName("id")
      if (self.playerId == null) {
        self.playerId = 0;
        self.byPlayer = false;
      } else {
        self.byPlayer = true;
      }
      count = self.refresh;
    },
    setupDefaults: function () {
port="3000"
      axios.defaults.baseURL = "https://localhost:" + port;
    },
    fetchPlayers: function () {
      const self = this;
      self.loading = true;
      request = "/getScores";
      axios
        .get(request)
        .then(function (response) {
          self.loading = false;
          const players = response.data;
          self.players = players;
        })
        .catch(onError.bind(self));
    },
    fetchPlayer: function () {
      const self = this;
      self.loading = true;
      request = "/player/" + self.playerId;
      axios
        .get(request)
        .then(function (response) {
          self.loading = false;
          const players = response.data;
          self.players = players;
        })
        .catch(onError.bind(self));
    },
    fetchStats: function () {
      const self = this;
      self.loading = true;
      request = "/stats";
      axios
        .get(request)
        .then(function (response) {
          self.loading = false;
          const stats = response.data;
          console.log(stats)
          self.stats = stats;
        })
        .catch(onError.bind(self));
    },
  },
});