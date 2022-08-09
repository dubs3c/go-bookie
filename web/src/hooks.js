
/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) { 
    const response = await resolve(event);

    // Should check here if auth store is marked as logged in

    return response;
}


/** @type {import('@sveltejs/kit').HandleError} */
export function handleError({ error, event }) {
    // example integration with https://sentry.io/
    // Sentry.captureException(error, { event });
}
