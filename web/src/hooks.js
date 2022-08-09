
/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) { 

    // this seems like such a hack...
    const cookies = event.request.headers.get("cookie");
    
    console.log(event.request.url);
    console.log(event.request.headers);
    
    const response = await resolve(event);
    
    if(!response.headers.get("cookie")) {
        response.headers.set("cookie", cookies);
    }
    console.log(response.url);
    console.log(response.headers);

    return response
}


/** @type {import('@sveltejs/kit').HandleError} */
export function handleError({ error, event }) {
    // example integration with https://sentry.io/
    // Sentry.captureException(error, { event });
}
