import { EventParams } from './event_params.tsx'
import {serverUrl} from "../main.tsx";

// event mapping, first level is event name, second level is component name, third level is corresponding params
export let events = new Map<string, Map<string, object>>()

export const EventName = {
    'click_event': 1,
    'page_view_event': 1,
}

export const DefaultEventSettingParams = {
    SAMPLE_FREQ: 1000,
    VIEW_DURATION: 5000,
}

export const { registerEvent } = {
    registerEvent: function (eventName, eventParams): boolean {
        // validate event name
        if (!(eventName in EventName))
        {
            console.log('Unknown track point event name')
            return false
        }

        // create event name
        if (!events.get(eventName))
        {
            events.set(eventName, new Map<string, object>())
        }

        // get event params
        for (const commonParam in EventParams)
        {
            if (!(commonParam in eventParams))
            {
                console.log('Unknown track point common params')
                return false
            }
        }

        for (const key in eventParams)
        {
            events.get(eventName)?.set(key, eventParams[key])
        }

        return true
    }
}

export const { sendEvent } = {
    sendEvent: function (eventName: string, params) {
        const eventStr = JSON.stringify({event_name: eventName, event_params: JSON.stringify(params)})
        const dataToSend = {
            'event':eventStr
        }

        console.log(dataToSend)

        const apiUrl = serverUrl + '/commerce/send-event'

        return fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(dataToSend)
        })
        .then(response => response.json())
        .then(result => {
            if (result['status_code'] == 0)
            {
                console.log('Send track point event success:', result)
            }
            else
            {
                return result['status_msg']
            }
        })
        .catch(error => {
            console.error('Send track point event error:', error);
            throw error
        });
    }
}

