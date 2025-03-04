import './event_dashboard.css'
import {serverUrl} from "../main.tsx";
import {useState} from "react";
import { Modal } from "react-bootstrap"
import {EventParams} from "./event_params.tsx";

function getPageViewNum(events) {
    return events.length
}

function getUniqueViewNum(events) {
    let cnt = new Map<string, null>
    for (const event of events)
    {
        cnt.set(event['username'], null)
    }

    return cnt.size
}

export function EventDashboard() {
    const [needQuery, setNeedQuery] = useState(true)
    const [events, setEvents] = useState([])

    function queryEvents(offset, limit) {
        const apiUrl = serverUrl + '/commerce/query-events'
        let params = JSON.stringify({offset: offset, limit: limit})

        // query events
        fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: params
        })
            .then(response => response.json())
            .then(result => {
                if (result['status_code'] != -1)
                {
                    setEvents(result['events'])
                    setNeedQuery(false)
                    return true
                }
                else
                {
                    console.log(result['status_msg'])
                    return false
                }
            })
            .catch(error => {
                console.error('Query track point event error:', error);
                throw error
            });
    }

    // prepare events data
    let eventNames = []
    if (needQuery)
    {
        queryEvents(0, 0)
    }
    else
    {
        // get unique event name
        let cnt = new Map<string, null>
        for (const event of events)
        {
            cnt.set(event['event_name'], null)
        }

        for (const key of cnt.keys())
        {
            eventNames.push(key)
        }
    }

    /* event name ui */
    const [selectedEventName, setSelectedEventName] = useState('')
    const [selectedEventId, setSelectedEventId] = useState(-1)
    const [selectedEventIds, setSelectedEventIds] = useState([])
    const [selectedEventParams, setSelectedEventParams] = useState([])
    function EventNameButton(key, eventName) {
        function onClick() {
            setSelectedEventName(eventName)

            let id = []
            let params = []
            for (const event of events)
            {
                if (event['event_name'] == selectedEventName)
                {
                    id.push(event['event_id'])
                    params.push(event['event_params'])
                }
            }
            setSelectedEventIds(id)
            setSelectedEventParams(params)

            let PV = getPageViewNum(params)
            let UV = getUniqueViewNum(params)
            setPrintInfo('PV: ' + PV + '\n UV:' + UV)
        }

        return (
            <button className='event_name_button' key={key} onClick={onClick}>
                {eventName}
            </button>
        )
    }

    /* event params ui */
    function EventParamUnit(key, event_id, param) {
        function onClick() {
            setSelectedEventId(event_id)
            setPrintInfo('Selected event id ' + event_id)
        }

        return (
            <div key={key} className='event_param_unit' onClick={onClick}>
                {key + '=' + param}
            </div>
        )
    }

    function clickDeletedButton() {
        const deleteEventUrl = serverUrl + '/commerce/delete-event'
        fetch(deleteEventUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({event_id: selectedEventId})
        })
            .then(response => response.json())
            .then(data => data['status_code'] != -1)
            .then(res => {
                if (res)
                {
                    setPrintInfo('Event ' + selectedEventId + ' deleted')
                    queryEvents(0, 0)
                }
                else
                {
                    setPrintInfo('Delete failed')
                }
            })
            .catch(error => {
                console.log(error)
            })
    }

    function EventParamsUnit(key, event_id, eventParams) {
        // parsing params to show in the unit
        let params = JSON.parse(eventParams)
        let paramList = []
        for (const idx in params)
        {
            paramList.push(EventParamUnit(idx, event_id, params[idx]))
        }

        return (
            <div key={key}>
                {paramList}
            </div>
        )
    }

    /* print board ui */
    const [printInfo, setPrintInfo] = useState('')
    function PrintBoard() {
        return (
            <div>
                {printInfo}
            </div>
        )
    }

    return (
        <div className='event_dashboard'>
            <div className='event_name_board'>
                {eventNames.map((eventName, idx) => {
                    return EventNameButton(idx, eventName)
                })}
            </div>
            <div className='event_params_board'>
                <div className='event_params_unit'>
                    {selectedEventParams.map((param, idx) => {
                        return EventParamsUnit(idx, selectedEventIds[idx], param)
                    })}
                </div>
                <div className='event_params_button'>
                    <button onClick={clickDeletedButton}>Delete</button>
                </div>
            </div>
            <div className='print_board'>
                <PrintBoard/>
            </div>
        </div>
    )
}
