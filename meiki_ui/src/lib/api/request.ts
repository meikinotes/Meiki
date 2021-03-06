export class StatusNotOkError extends Error {
    constructor(msg: string) {
        super(msg)
    }
}

async function ensureStatusOK(response: Response) {
    if (response.status !== 200) throw new StatusNotOkError(await response.json())
}

export async function makeRequest(url: string, method: string, body: any = {}) {
    // throws StatusNotOkError of the request was successfully made but the returned response had a non 200 status code
    // StatusNotOkError will have the response body
    // throws TypeError if it was not able to make the request to the URL
    const baseUrl = import.meta.env.VITE_MEIKI_SERVER_URL || "http://localhost:8080"
    const absUrl = new URL(url, baseUrl)

    const username = localStorage.getItem("username")
    const token = localStorage.getItem("token")

    let requestOptions: any = {
        method,
        headers: {
            "Content-Type": "application/json",
            "X-Username": username,
            "X-Token": token,
        },
        credentials: "include",
    }

    if (["POST", "PUT"].includes(method)) {
        requestOptions.body = JSON.stringify(body)
    }

    const response = await fetch(absUrl.toString(), requestOptions)
    await ensureStatusOK(response)
    return response.json()
}

export function formatRequestError(err: any, action: string) {
    if (err instanceof StatusNotOkError) return err.message

    //TODO: Make switch statement with Enums
    return `An error has occurred while ${action}, unable to connect to server`
}
