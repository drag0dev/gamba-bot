import React from 'react'

const Info = () => {

    return(
        <div className='info'>
            <h2>
               Currently supported: 
            </h2>

            <div className='keydrop'>
                <p>
                    Keydrop
                </p>
                <img src='keydrop.png'/>
            </div>

            <div className='csgocases'>
                <p>
                    CSGOCases
                </p>

                <img src='csgocases.jpg' />
                
            </div>

        </div>
    );
}

export default Info;