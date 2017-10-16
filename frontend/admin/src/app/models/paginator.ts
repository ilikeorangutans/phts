export class Paginator {

    readonly direction: string;
    readonly count: number;
    readonly paginatorType: PaginatorType;

    public static newTimestampPaginator(column: string): Paginator {
        return new Paginator(
            new TimestampPaginatorType(
                new Date(),
                column
            )
        );
    }

    constructor(
        theType: PaginatorType,
        theCount: number = 12,
        theDirection: string = 'desc'
    ) {
        this.paginatorType = theType;
        this.count = theCount;
        this.direction = theDirection;
    }

    toQueryString(): string {
        let data = new Array<string>();
        data.push(`direction=${this.direction}`);
        data.push(`count=${this.count}`);

        data = data.concat(this.paginatorType.addColumnAndValue());

        return data.join('&');
    }
}

export interface PaginatorType {
    addColumnAndValue(): Array<string>;
}

export class TimestampPaginatorType implements PaginatorType {

    readonly column: string;
    readonly timestamp: Date;

    constructor(theTimestamp: Date, theColumn: string) {
        this.column = theColumn;
        this.timestamp = theTimestamp;
    }

    addColumnAndValue(): Array<string> {
       return [
           `column=${this.column}`,
           `prevTimestamp=${this.timestamp.toISOString()}`
       ];
    }
}
