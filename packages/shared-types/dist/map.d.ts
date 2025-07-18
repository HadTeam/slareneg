import { Block, Owner } from './block';
export interface Size {
    width: number;
    height: number;
}
export interface Pos {
    x: number;
    y: number;
}
export interface Info {
    id: string;
    name: string;
    desc: string;
}
export type Sight = boolean[][];
export type Blocks = Block[][];
export interface Map {
    isEmpty(): boolean;
    block(pos: Pos): Block | null;
    blocks(): Blocks;
    setBlock(pos: Pos, b: Block): void;
    setBlocks(blocks: Blocks): void;
    size(): Size;
    info(): Info;
    roundStart(roundNum: number): void;
    roundEnd(roundNum: number): void;
    fog(owner: Owner[], sight: Sight): void;
}
export declare function sizeToString(s: Size): string;
export declare function posToString(p: Pos): string;
export declare function isPosValid(s: Size, p: Pos): boolean;
export declare function infoToString(i: Info): string;
