export interface IPhishingEvent {
    id: number;
    name: string,
    createdAt: Date,
    brand: string,
    description: string,
    maliciousUrl: string,
    domainRegistrationDate: Date,
    keyword: string[],
    status: "todo" | "in progress" | "done",
    comments: IComment[],
    dnsRecords: string[]
}

export interface IComment {
    createdAt: Date,
    updatedAt: Date,
    comment: string,
    username: string
}

export interface IDnsRecords{
    ns: string,
    a: string
    mx: string
}

export interface PhishingEventSearchData {
    name?: string,
    startDate?: Date,
    endDate?: Date,
    brand?: string,
    domainName?: string,
    keywords?: Array<string>,
}
