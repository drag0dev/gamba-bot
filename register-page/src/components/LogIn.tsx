import React from 'react'
import DiscordInfo from '../DiscordInfo';

const LogIn = () => {

    const onClickSubscribe = () => {}

    let url = `https://discord.com/api/oauth2/authorize?response_type=token&client_id=${DiscordInfo.ClientID}&scope=identify`

    return(
        <div className='login'>

            <p>
                Want to subscribe to Gamba Bot?
            </p>

            <br />

            <p className='clickme' onClick={onClickSubscribe}>
                <a rel='noopener noreferrer' href={url}>
                    Click me
                </a>
            </p>

        </div>
    );
}

export default LogIn;