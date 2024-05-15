import {whoisInstance} from "@api/config/axios.ts";

export const getDnsRecords = async (url: string) => {

    const data = await whoisInstance.get("/whoisserver/WhoisService", {
        params: {
            outputFormat: "json",
            apiKey: "at_X6tVoNL7IJMGV6CIPmKCVI2DYdIRQ",
            domainName: url
        },
    })
    return data.data

}
