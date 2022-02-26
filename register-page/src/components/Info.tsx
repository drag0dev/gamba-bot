import React from 'react'

const Info = () => {

    return(
        <div className='info'>
            
            <div>

                <h3>
                    Currently supported: 
                </h3>

            </div>

            <div className='keydrop'>

                <div>
                    <p>
                        Keydrop
                    </p>
                </div>

                <div>
                    <img src='keydrop.png'/>
                </div>

            </div>

            <div className='csgocases'>

                <div>
                    <p>
                        CSGOCases
                    </p>
                </div>

                <div>
                    <img src='csgocases.jpg' />
                </div>
                
            </div>

        </div>
    );
}

export default Info;