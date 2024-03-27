const express = require('express')
const passport = require('passport')
const session = require('express-session')
const SteamStrategy = require('passport-steam').Strategy;

const {axios, instance} = require('./configure_https_agent');

// TODO: THIS IS STILL A CRUDE IMPLEMENTATION OF 
// STEAM ACCOUNT LINK

const SYNC_ENDPOINT_HOST = 'localhost:3000/v1/games'
const LINK_ENDPOINT_HOST = `localhost:3000/v1/account`



// TODO: Handle certificate
const httpsAgent = new https.Agent({
  cert: fs.readFileSync(CERT_FILE),
  key: fs.readFileSync(KEY_FILE),
  ca: fs.readFileSync(ROOT_CA),
});


async function getSteamOwnedGames(steamID) {
  url = 'https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/'

  //https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?steamid=76561198221895016&include_played_free_games=1&include_appinfo=1&format=json
  const res = await axios.get(
    url, {
      params: {
        steamid: steamID,
        include_played_free_games: 1,
        include_appinfo: 1,
        format: 'json',
        key: process.env.API_KEY
      }
    }
  )

  
  return res.data.response
}


// Passport session setup.
//   To support persistent login sessions, Passport needs to be able to
//   serialize users into and deserialize users out of the session.  Typically,
//   this will be as simple as storing the user ID when serializing, and finding
//   the user by ID when deserializing.  However, since this example does not
//   have a database of user records, the complete Steam profile is serialized
//   and deserialized.
// Passport session setup.
//   To support persistent login sessions, Passport needs to be able to
//   serialize users into and deserialize users out of the session.  Typically,
//   this will be as simple as storing the user ID when serializing, and finding
//   the user by ID when deserializing.  However, since this example does not
//   have a database of user records, the complete Steam profile is serialized
//   and deserialized.
passport.serializeUser(function(user, done) {
  done(null, user);
});

passport.deserializeUser(function(obj, done) {
  done(null, obj);
});

// Use the SteamStrategy within Passport.
//   Strategies in passport require a `validate` function, which accept
//   credentials (in this case, an OpenID identifier and profile), and invoke a
//   callback with a user object.


var app = express();

// configure Express
app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');

app.use(session({
    secret: 'your secret',
    name: 'name of session id',
    resave: true,
    saveUninitialized: true}));

// Initialize Passport!  Also use passport.session() middleware, to support
// persistent login sessions (recommended).
// app.use(passport.initialize());
app.use(passport.session());
app.use(express.static(__dirname + '/../../public'));

// GET /auth/steam
//   Use passport.authenticate() as route middleware to authenticate the
//   request.  The first step in Steam authentication will involve redirecting
//   the user to steamcommunity.com.  After authenticating, Steam will redirect the
//   user back to this application at /auth/steam/return

const configurePassport = (req, res, next) => {
  const { username } = req.params;
  const { failRedirect, successRedirect } = req.query;

  passport.use(new SteamStrategy({
    returnURL: `http://localhost:3000/middleware/steam/return/${username}?failureRedirect=${failRedirect}&successRedirect=${successRedirect}`,
    realm: 'http://localhost:3000/',
    profile: false,
  },
  function(identifier, profile, done) {
    // asynchronous verification, for effect...
    process.nextTick(function () {
  
      // To keep the example simple, the user's Steam profile is returned to
      // represent the logged-in user.  In a typical application, you would want
      // to associate the Steam account with a user record in your database,
      // and return that user instead.
      profile.identifier = identifier;
      return done(null, profile);
    });
  }
  ));
  next()
}
app.use('/middleware/steam/:username', configurePassport);
app.use('/middleware/steam/return/:username', configurePassport);

app.get('/middleware/steam/:username',
  passport.authenticate('steam', { failureRedirect: '/' }),
  function(req, res) {
    res.redirect('/');
  });

// GET /auth/steam/return
//   Use passport.authenticate() as route middleware to authenticate the
//   request.  If authentication fails, the user will be redirected back to the
//   login page.  Otherwise, the primary route function function will be called,
//   which, in this example, will redirect the user to the home page.
app.get('/middleware/steam/return/:username',
  passport.authenticate('steam', { failureRedirect: '/' }),
  async function(req, res) {
    console.log(req.params.username);

    const { failRedirect, successRedirect } = req.query;

    splitStr = req.user.identifier.split("/");
    id = splitStr[splitStr.length - 1];

    const userGames = await getSteamOwnedGames(id);
    
    // filter data in userGames
    filteredUserGames = userGames.games.map(obj => {
      return {
        name: obj.name,
        app_id:obj.appid,
        icon_url:`"http://media.steampowered.com/steamcommunity/public/images/apps/${obj.appid}/${obj.img_icon_url}.jpg"`
      }
    });


    data = {
      games: filteredUserGames
    };

    console.log(data);

    // link steam account
    let resp = instance.post(`https://${LINK_ENDPOINT_HOST}/${username}/steam`, 
    {
      steamid: id
    }, 
    {
      headers: {
        'Content-Type':'application/json',
      }
    });

    console.log(resp.data.response);
    
    if (resp.status != 200) {
      // res.redirect(failRedirect);
      return;
    }


    // sync games
    resp = await instance.post(`https://${SYNC_ENDPOINT_HOST}/${username}/sync`, 
    {
      games: data,
    }, 
    {
      headers: {
        'Content-Type':'application/json',
      }
    }); 

    if (resp.status != 200) {
      resp = await instance.delete(`https://${LINK_ENDPOINT_HOST}/${username}/steam`);
      // res.redirect(failRedirect);
    }


    const res = instance.post('')
    // res.redirect(successRedirect);
});

  

app.listen(3000);
