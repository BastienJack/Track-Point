import Product from '../assets/product.jpg'
import './item.css'
import {useEffect, useRef, useState} from "react";
import {sendEvent} from '../track-point/event.tsx'
import { registerEvent } from "../track-point/event.tsx";
import { EventParams } from "../track-point/event_params.tsx";

function ItemNavigator() {
    function RedirectToEventDashboard() {
        window.location.replace('event-dashboard')
    }

    return (
        <div className='item_navigator'>
            <button onClick={RedirectToEventDashboard}>
                Event Dashboard
            </button>
        </div>
    )
}

const rowIncNum = 4
function Item(idx) {
    let componentName = 'item' + idx
    let eventParams = EventParams

    useEffect(() => {
        // register event
        const username = localStorage.getItem('username')
        eventParams = {
            'username': username,
            'component_name': componentName,
        }

        if (!registerEvent('click_event', eventParams))
        {
            console.log(componentName + ' click_event register failed')
        }

        return () => {
        }
    }, []);

    function itemClickEvent() {
        // check login
        if (localStorage.getItem('username') == null)
        {
            alert("Please login first")
            window.location.replace("login")
            return
        }

        window.open("item-detail?item_name=" + componentName)
        sendEvent('click_event', eventParams)
    }

    return (
        <div id={componentName} key={idx} className='item' onClick={itemClickEvent}>
            <img src={Product} alt={'item'}/>
            <span>Item{idx}</span>
        </div>
    )
}

function ItemGrid() {
    const [itemNum, setItemNum] = useState(16)

    let items = []
    for (let i = 0; i < itemNum; i++)
    {
        items.push(Item(i))
    }

    const gridRef = useRef(null)
    const [isBottom, setIsBottom] = useState(false)

    useEffect(() => {
        const handleScroll = () => {
            /* implement scroll to end listener */
        }

        if (isBottom)
        {
            setItemNum(itemNum + rowIncNum)
        }

        window.addEventListener('scroll', handleScroll)

        return () => {
            window.removeEventListener('scroll', handleScroll)
        }
    }, [])

    return (
        <div className='item_grid'  ref={gridRef}>
            {items}
        </div>
    )
}

export default function ItemPage() {
    return (
        <div className='item_page'>
            <ItemNavigator/>
            <ItemGrid/>
        </div>
    )
}
