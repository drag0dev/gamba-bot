import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

const Header = () => {
    const navigate = useNavigate();
    const location = useLocation();

    return(
        <div className='header'>
                <h1>
                    Gamba bot
                </h1>


            {location.pathname.includes('/redirect') &&

                <p onClick={() => {navigate('/')}} className='home-click'>
                    Home
                </p>
            }

        </div>       
    );
}

export default Header;