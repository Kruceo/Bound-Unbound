import React, { createContext, useContext, ReactNode, useState } from 'react';
import BeautyBox from '../components/BeautyBox';
import "./NotificationContext.less";
import Ico from '../components/Ico';

// Define the shape of our context
interface NotificationContextType {
    spawnNotification: (text: string) => void
}

export const NotificationContext = createContext<NotificationContextType>({ spawnNotification: () => null });


interface NotificationProviderProps {
    children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
    const [currentNotification, setCurrentNotification] = useState<ReactNode>(null);

    function spawnNotf(text: string) {
        setCurrentNotification(<><Notification text={text} title='Notification'/></>)
        setTimeout(()=>setCurrentNotification(<></>),5000)
    }

    return (
        <NotificationContext.Provider value={{ spawnNotification: spawnNotf }}>
            {currentNotification}
            {children}
        </NotificationContext.Provider>
    );
};

function Notification(props: { title: string, text: string }) {
    return <>
        <BeautyBox className='notification'>
            <div className='image'>
                <Ico>error</Ico>
            </div>
            <div className='text'>
                <h2>{props.title}</h2>
                <p>{props.text}</p>
            </div>
            <div>
                {/* <button>Close</button> */}
            </div>
        </BeautyBox>
    </>
}