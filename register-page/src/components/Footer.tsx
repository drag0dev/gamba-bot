import React from 'react'

const Footer = () =>{

    return(
        <div className='footer'>
            
            <p>
                &copy; 2022 made by drag0dev 
            </p>

            <a target='_blank' rel='noopener noreferrer' href='https://www.github.com/drag0dev'>
                <img className='footer-img' src='github.png'></img>
                
            </a>
        
        </div>
    );
}

export default Footer;