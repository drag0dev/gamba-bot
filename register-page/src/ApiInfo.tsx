const ApiInfo = {
   EXISTS: process.env.REACT_APP_API_EXISTS || 'https://gamba-bot-api.herokuapp.com/exists',
   SUBSCRIBE: process.env.REACT_APP_API_SUBSCRIBE || 'https://gamba-bot-api.herokuapp.com/subscribe',
   UNSUBSCRIBE: process.env.REACT_APP_API_UNSUBSCRIBE || 'https://gamba-bot-api.herokuapp.com/unsubscribe' 
}

export default ApiInfo;