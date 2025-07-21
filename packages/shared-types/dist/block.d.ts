export type Num = number;
export type Owner = number;
export type Name = string;
export interface Meta {
    name: Name;
    description: string;
}
export interface AllowMove {
    from: boolean;
    to: boolean;
    reason: string;
}
export interface Block {
    num(): Num;
    owner(): Owner;
    roundStart(roundNum: number): void;
    roundEnd(roundNum: number): void;
    allowMove(): AllowMove;
    moveFrom(num: Num): Num;
    moveTo(num: Num, owner: Owner): Block;
    fog(isOwner: boolean, isSight: boolean): Block;
    meta(): Meta;
}
