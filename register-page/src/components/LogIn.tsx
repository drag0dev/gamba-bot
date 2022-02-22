import React from 'react'

const LogIn = () => {

    const onClickSubscribe = () => {}

    return(
        <div className='login'>

            <p>
                Want to subscribe to Gamba Bot?
            </p>

            <br />

            <p className='clickme' onClick={onClickSubscribe}>
                <a rel='noopener noreferrer' href='https://discord.com/api/oauth2/authorize?response_type=token&client_id=942833925431119873&scope=identify'>
                    Click me
                </a>
            </p>

        </div>
    );
}

export default LogIn;