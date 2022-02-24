import React from 'react'
import DiscordInfo from '../DiscordInfo';

const LogIn = () => {
    let url = `https://discord.com/api/oauth2/authorize?response_type=token&client_id=${DiscordInfo.ClientID}&scope=identify`;
    let invite = 'https://discord.gg/BX2NzeG86r';
    
    return(
        <div className='login'>

            <div>

                <p>
                    Want to subscribe/unsubscribe to Gamba Bot?
                </p>

                <p className='clickme'>
                    <a rel='noopener noreferrer' href={url}>
                        Click me
                    </a>
                </p>

            </div>

            <div>

                <p>
                    Already subscribed? Join server to receive promocodes!
                </p>

                <p className='clickme'>
                    <a href={invite} rel='nooopener noreferrer' target='_blank'>
                        Join 
                    </a>
                </p>

            </div>

        
        </div>
    );
}

export default LogIn;