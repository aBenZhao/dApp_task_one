// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 计数器合约
contract Counter {
    uint256 public count;

    constructor() {
        count = 0;
    }

    function increment() external {
        count += 1; // 等价于 count = count + 1
    }
}