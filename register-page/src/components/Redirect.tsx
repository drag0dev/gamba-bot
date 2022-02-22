import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';

interface userInfo {
    id: string,
    username: string,
    discriminator: string,
    avatar: string
}

const Redirect = () =>{
    const location = useLocation();

    const getUserInfo = async (data: string[]) => {
        let res = await fetch('https://discord.com/api/v8/users/@me', {
            headers: {
                'Authorization': `${data[0]} ${data[1]}`,
            }

        })

        let userData: userInfo = await res.json()
        console.log(userData.id);
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

            <p>
                Processing your info...
            </p>

        </div>
    );
}

export default Redirect;