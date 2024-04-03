const express = require('express')
const axios = require('axios')
const passport = require('passport')
const session = require('express-session')
const SteamStrategy = require('passport-steam').Strategy;

const {instance} = require('./configure_https_agent');

// TODO: THIS IS STILL A CRUDE IMPLEMENTATION OF 
// STEAM ACCOUNT LINK

const SYNC_ENDPOINT_HOST = 'proxy:3000/v1/games'
const LINK_ENDPOINT_HOST = `proxy:3000/v1/account`

const port = process.env.port || 9000
const host = process.env.host || 'localhost'

if (process.env.STEAM_API_KEY == '') {
  console.log('Steam Api Key env var is not set. Exiting...')
  process.exit(1)
}

function checkAxiosError(error) {
  // Ref: https://axios-http.com/docs/handling_errors
  if (error.response) {
    // Response not on the type 2xx
    console.log(error.response.data)
    console.log(error.response.status)
  } else if (error.request) {
    // Request was amde but no response was received
    console.log(error.request)
  } else {
    // something bad happened when setting up the request
    console.log(`Error: ${error.message}`)
  }
}

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
        key: process.env.STEAM_API_KEY
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
    returnURL: `http://${host}:${port}/middleware/steam/return/${username}?failureRedirect=${failRedirect}&successRedirect=${successRedirect}`,
    realm: `http://${host}:${port}/`,
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

    const {username} = req.params

    const { failRedirect, successRedirect } = req.query;

    splitStr = req.user.identifier.split("/");
    id = splitStr[splitStr.length - 1];

    // link steam account
    try {
      resp = await instance.post(`http://${LINK_ENDPOINT_HOST}/${username}/steam`, 
        {
          steamid: id
        }, 
        {
          headers: {
            'Content-Type':'application/json',
          }
        });

      console.log(resp.data);
    } catch (error) {
      checkAxiosError(error);
      return;
    }

    try {
      const userGames = await getSteamOwnedGames(id);

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

      resp = await instance.post(`http://${SYNC_ENDPOINT_HOST}/${username}/sync`, 
      data, 
      {
        headers: {
          'Content-Type':'application/json',
        }
      }); 

      console.log(resp.data)

    } catch (error) {
      // TODO: redirect url
      checkAxiosError(error);
      resp = await instance.delete(`http://${LINK_ENDPOINT_HOST}/${username}/steam`)
      return;
    }
    

    return;
});

  
app.listen(port, '0.0.0.0');
console.log(`This service is running on ${host}`)
console.log(`Middleware accept connection on port 0.0.0.0:${port}`)
