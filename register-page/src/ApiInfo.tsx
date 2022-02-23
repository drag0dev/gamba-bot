const ApiInfo = {
   EXISTS: process.env.REACT_APP_API_EXISTS || 'http://localhost:8080/exists',
   SUBSCRIBE: process.env.REACT_APP_API_SUBSCRIBE || 'http://localhost:8080/subscribe',
   UNSUBSCRIBE: process.env.REACT_APP_API_UNSUBSCRIBE || 'http://localhost:8080/unsubscribe' 
}

export default ApiInfo;