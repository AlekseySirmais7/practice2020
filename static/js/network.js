import CurrentUser from "./currentUser";

export default class FetchModule {
    /**
     * Post request with json
     *
     * @param {string} url - request url
     * @param {map} parameters - parameters for fetch
     * @return {promise}
     */
    static fetchRequest({url,
                            method = 'post',
                            body = null,
                            headers = {
                                'Content-type': 'application/json; charset=UTF-8',
                                'X-CSRF-Token': CurrentUser.Data.token,
                            },
                        } = {}) {
        const jsonData = JSON.stringify(body);
        const options = {
            method: method,
            credentials: 'include',
            headers: headers,
            mode: 'cors',
            body: body == null ? null : jsonData,
        };
        return fetch(url, options);
    };
}