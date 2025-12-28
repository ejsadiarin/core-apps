const url = "http://localhost:8080"

export async function fetchSystemStatus() {
    const res = await fetch(url + "/api/system/stats")
    if (!res.ok) {
        throw new Error(`Error fetching ${res.status}`)
    }
    return res.json()
}

export async function fetchServicesStats() {
    const res = await fetch(url + "/api/services")
    if (!res.ok) {
        throw new Error(`Error fetching ${res.status}`)
    }
    return res.json()
}
