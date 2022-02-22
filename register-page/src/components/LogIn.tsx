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
                Click me
            </p>

        </div>
    );
}

export default LogIn;