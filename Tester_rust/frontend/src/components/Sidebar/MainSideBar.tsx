import React, { useRef, useState } from 'react'
import { SidebarClasses } from '../../App'
import "./MainSidebar.css"



interface MainSideBarProps{
}

const MainSideBar = ({}:MainSideBarProps) => {
     const [isVisible, setIsVisible] = useState(false)
    const hoverTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

    const handleMouseEnter = () => {
        hoverTimer.current = setTimeout(() => {
            setIsVisible(true)
        }, 150)
    }

    const handleMouseLeave = () => {
        if (hoverTimer.current) {
            clearTimeout(hoverTimer.current)
            hoverTimer.current = null
        }

        setIsVisible(false)
    }

    return (
        <>
            <div className={SidebarClasses.trigger} onMouseEnter={handleMouseEnter}/>

            {isVisible && (
                <div className={SidebarClasses.container} onMouseLeave={handleMouseLeave}>
                    <div className='sidebar-first'>
                        <button className={SidebarClasses.button}>abc</button>
                        <button className={SidebarClasses.button}>abc</button>
                    </div>

                    <hr />

                    <div className='sidebar-second'>
                        <button className={SidebarClasses.button}>abc</button>
                        <button className={SidebarClasses.button}>abc</button>
                    </div>

                    <hr />
                    
                    <div className='sidebar-third'>
                        <button className={SidebarClasses.button}>abc</button>
                        <button className={SidebarClasses.button}>abc</button>
                    </div>
                    
                  
                </div>
            )}
        </>
    )
}

export default MainSideBar