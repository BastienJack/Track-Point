import {serverUrl} from "../main.tsx";

export let EventParams = {}

export function readEventCommonParams() {
    const getParamsUrl = serverUrl + '/commerce/get-common-params'

    fetch(getParamsUrl)
       .then(response => response.json())
       .then(data => {
           EventParams = data
           console.log(EventParams)
       })
        .catch(error => {
            console.log(error)
        })
}

export function addEventCommonParams(key: string, defaultValue: string) {
    const addParamsUrl = serverUrl + '/commerce/add-common-params'

    fetch(addParamsUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({'key': key, 'value': defaultValue})
    })
        .then(response => {
            console.log(response.json())
        })
        .then(data => {
            if (data['status_code'] != -1)
            {
                EventParams[key] = defaultValue
            }
        })
        .catch(error => {
            console.log(error)
        })
}
