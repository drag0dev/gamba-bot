import React, { useState, useEffect } from 'react';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import ApiInfo from '../ApiInfo'

interface userInfo {
    id: string,
    username: string,
    discriminator: string,
    avatar: string
}

const Redirect = () =>{
    const location = useLocation();
    const [message, setMessage] = useState('Processing your info...');
    const [userLoading, setUserLoading] = useState(true);
    const [userState, setUserState] = useState(false);
    const [userData, setUserData] = useState <userInfo> ({
        id: '',
        username: '',
        discriminator: '',
        avatar: ''
    });
    const navigate = useNavigate();

    const getUserInfo = async (data: string[]) => {
        let res = await fetch('https://discord.com/api/v8/users/@me', {
            headers: {
                'Authorization': `${data[0]} ${data[1]}`,
            }

        })

        let tempUserData: userInfo = await res.json();
        await setUserData(tempUserData);

        // check if user is subscribed
        res = await fetch(ApiInfo.EXISTS, {
            method: 'POST',
            headers: {
                'Content-type': 'application/json'
            },
            body: JSON.stringify({id: tempUserData.id})
        })

        let resData = await res.json()
        if (res.status == 200){
                setMessage(`Hello, ${tempUserData.username!}`)
                if (resData.found == 'true')setUserState(true)
                setUserLoading(false);
        }

        else{
            setMessage('There was a problem processing your request, please try again!')
            const timer = setTimeout(()=>{
                navigate('/')
            }, 2000);
        }
    }

    const onClickSubscribe = async () => {
        setUserLoading(true)
        let res = await fetch(ApiInfo.SUBSCRIBE, {
            method: 'POST',
            headers: {
                'Content-type': 'application/json'
            },
            body: JSON.stringify({id: userData.id})
        });

        if (res.status == 200){
            setMessage('Successfully subscribed, returning home... (DON\'T FORGET TO JOIN THE SERVER)')
            const timer = setTimeout(()=>{
                navigate('/')
            }, 2000);
        }
        else{
            setMessage('There was a problem processing your request, please try again!')
            const timer = setTimeout(()=>{
                navigate('/')
            }, 2000);
        }
    }

    const onClickUnsubscribe = async () => {
        setUserLoading(true)
        let res = await fetch(ApiInfo.UNSUBSCRIBE, {
            method: 'POST',
            headers: {
                'Content-type': 'application/json'
            },
            body: JSON.stringify({id: userData.id})
        });

        if (res.status == 200){
            setMessage('Successfully unsubscribed, returning home...')
            const timer = setTimeout(() => {
                navigate('/')
            }, 2000);
        }
        else{
            setMessage('There was a problem processing your request, please try again!')
            const timer = setTimeout(()=>{
                navigate('/')
            }, 2000);
        }

    }

    useEffect(()=>{
        let parsedParams: string[] = []; // token_type, access_token, expires_in, scope

        if (location){
            let urlParams = location.hash
            let params = urlParams.replaceAll('?', '').split('&')

            params.forEach((param, index)=>{
                let temp = param.split('=')
                parsedParams[index] = temp[1]
            });
            getUserInfo(parsedParams);
        }
    }, [location]);

    return(
        <div className='redirect'>
            <div className='redirect-info'>

                <div>
                    <p>
                      {message}
                    </p>
                </div>

                {!userLoading &&
                    <div>
                        {userState
                        ? <button onClick={onClickUnsubscribe}>Unsubscribe</button>
                        : <button onClick={onClickSubscribe}>Subscribe</button>
                        }
                    </div>
                }

            </div>

        </div>
    );
}

export default Redirect;
