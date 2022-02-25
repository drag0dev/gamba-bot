import React from 'react';
import { useNavigate, useLocation, Link} from 'react-router-dom';

const Header = () => {
    const navigate = useNavigate();
    const location = useLocation();

    return(
        <div className='header'>
                <h1>
                    Gamba bot
                </h1>


            <div className='options'>

                <p onClick={() => {navigate('/')}} >
                    Home
                </p>

                <p >
                    <Link to='/commands'>Commands</Link>
                </p>

            </div>

        </div>       
    );
}

export default Header;