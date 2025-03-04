import Product from "../assets/product.jpg";
import {useEffect, useState} from "react";
import {DefaultEventSettingParams, registerEvent, sendEvent} from "../track-point/event.tsx";
import {EventParams} from "../track-point/event_params.tsx";


export default function ItemDetail() {
    const url = new URL(window.location.href)
    const item_name = url.searchParams.get('item_name')!

    let eventParams = EventParams
    const [isViewed, setIsViewed] = useState(false)
    useEffect(() => {
        // register event
        const username = localStorage.getItem('username')
        eventParams = {
            'username': username,
            'component_name': item_name,
        }

        if (!registerEvent('page_view_event', eventParams))
        {
            console.log(item_name + ' page_view_event register failed')
        }

        // send page view event
        const timer = setTimeout(() => {
            // check login
            if (localStorage.getItem('username') == null)
            {
                alert("Please login first")
                window.location.replace("login")
                return
            }

            if (!isViewed)
            {
                setIsViewed(true)
                sendEvent('page_view_event', eventParams)
            }
        }, DefaultEventSettingParams.VIEW_DURATION)

        return () => {
            clearTimeout(timer)
        }
    }, []);

    return (
        <div>
            <img src={Product} alt={item_name}/>
        </div>
    )
}
