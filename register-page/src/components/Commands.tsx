import React from 'react';

const Commands = () =>{

    return(
        <div className='commands'>

            <div className='commands-header'>
                <h1>
                    Commands
                </h1>
            </div>

            <div className='command'>

                <div className='option-name'>
                    <h3>!subscribe</h3>
                </div>

                <div>
                        <p>
                            Command subscribes you to the bot, needs to be entered in a server where Gamba Bot is a member.
                        </p>
                </div>

            </div>

            <div className='command'>

                <div className='option-name'>
                    <h3>!unsubscribe</h3>
                </div>

                <div>
                        <p>
                            Command unsubscribes you from the bot, needs to be entered in a server where Gamba Bot is a member.
                        </p>
                </div>

            </div>

            <div className='command'>

                <div className='option-name'>
                    <h3>!bind</h3>
                </div>

                <div>
                        <p>
                            Command that binds bot to a channel that the command has been entered, all new codes will be sent to that channel. Needs to be entered in a server where Gamba Bot is a member.
                        </p>
                </div>

            </div>

            <div className='command'>

                <div className='option-name'>
                    <h3>!unbind</h3>
                </div>

                <div>
                        <p>
                            Command that unbinds bot form a channel, needs to be entered in a server where Gamba Bot is a member.
                        </p>
                </div>

            </div>
        
        </div>
    );
}

export default Commands;