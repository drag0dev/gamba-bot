import React from 'react';

const Commands = () =>{

    return(
        <div className='commands'>

            <div className='commands-header'>
                <h1>
                    Commands
                </h1>
            </div>

            <div className='small-header'>
                <h2 className='permissions-h'>
                    permissions
                </h2>

                <h2 className='command-h'>
                    command
                </h2>

                <h2 className='description-h'>
                    description
                </h2>

            </div>

            <div className='command'>

                <div className='perm'>
                    <p>
                        none 
                    </p>
                </div>

                <div className='option-name'>
                    <h3>$subscribe</h3>
                </div>

                <div className='command-desc'>
                    <p>
                        subscribe to the bot 
                    </p>
                </div>

            </div>

            <div className='command'>

                <div className='perm'>
                    <p>
                        none 
                    </p>
                </div>

                <div className='option-name'>
                    <h3>$unsubscribe</h3>
                </div>

                <div className='command-desc'>
                        <p>
                            unsubscribe from the bot 
                        </p>
                </div>

            </div>

            <div className='command'>
                
                <div className='perm'>
                    <p>
                        admin or manage server 
                    </p>
                </div>

                <div className='option-name'>
                    <h3>$bind</h3>
                </div>

                <div className='command-desc'>
                        <p>
                            bind bot to a channel that the command has been entered in, all future codes will be sent to that channel
                        </p>
                </div>

            </div>

            <div className='command'>

                <div className='perm'>
                    <p>
                        admin or manage server 
                    </p>
                </div>

                <div className='option-name'>
                    <h3>$unbind</h3>
                </div>

                <div className='command-desc'>
                        <p>
                            unbind bot from a channel
                        </p>
                </div>

            </div>
            
            <div className='command'>

                <div className='perm'>
                    <p>
                        none 
                    </p>
                </div>

                <div className='option-name'>
                    <h3>$grab csgocases/keydrop</h3>
                </div>

                <div className='command-desc'>
                        <p>
                            get last 5 codes for a specified website in your DMs
                        </p>
                </div>

            </div>
        
        </div>
    );
}

export default Commands;