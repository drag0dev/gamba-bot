import React from 'react'
import DiscordInfo from '../DiscordInfo';

const LogIn = () => {
    let url = `https://discord.com/api/oauth2/authorize?response_type=token&client_id=${DiscordInfo.ClientID}&scope=identify`;
    let invite = 'https://discord.com/invite/bWQQnC9CCe';
    let botInvite = `https://discord.com/api/oauth2/authorize?client_id=942833925431119873&permissions=8&scope=applications.commands%20bot`;

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

            <div>

                <p>
                    Want Gamba Bot on your server?
                </p>

                <p className='clickme'>
                    <a  rel='noopener noreferrer' target='_blank' href={botInvite}>
                        Invite Gamba Bot
                    </a>
                </p>

            </div>


        </div>
    );
}

export default LogIn;
