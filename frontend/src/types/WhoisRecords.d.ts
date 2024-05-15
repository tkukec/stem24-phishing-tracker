export interface IDomainRecords {
    createdDate: Date,
    updatedDate: Date,
    expiresDate: Date,
}

interface IWhoisRecords {
    WhoisRecord: IDomainRecords
}
