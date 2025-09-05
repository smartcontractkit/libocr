// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

library LowLevelCallLib {

  uint256 private constant CALL_WITH_EXACT_GAS_CUSHION = 5_000;

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function callWithExactGasEvenIfTargetIsNoContract(
    uint256 _gasAmount,
    address _target,
    bytes memory _data
  )
    external
    returns (bool sufficientGas)
  {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      let g := gas()
      // Compute g -= CALL_WITH_EXACT_GAS_CUSHION and check for underflow. We
      // need the cushion since the logic following the above call to gas also
      // costs gas which we cannot account for exactly. So cushion is a
      // conservative upper bound for the cost of this logic.
      if iszero(lt(g, CALL_WITH_EXACT_GAS_CUSHION)) {
        g := sub(g, CALL_WITH_EXACT_GAS_CUSHION)
        // If g - g//64 <= _gasAmount, we don't have enough gas. (We subtract g//64
        // because of EIP-150.)
        if gt(sub(g, div(g, 64)), _gasAmount) {
          // Call and ignore success/return data. Note that we did not check
          // whether a contract actually exists at the _target address.
          pop(call(_gasAmount, _target, 0, add(_data, 0x20), mload(_data), 0, 0))
          sufficientGas := true
        }
      }
    }
  }

}
