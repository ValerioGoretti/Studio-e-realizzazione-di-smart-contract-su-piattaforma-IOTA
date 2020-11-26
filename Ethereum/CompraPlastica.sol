pragma solidity >=0.4.21 <0.8.0;

import './TokenPlastica.sol';

contract CompraPlastica
{
    TokenPlastica tokenPlasticaProvider=TokenPlastica(address(0x5FD6eB55D12E759a21C09eF703fe0CBa1DC9d88D));
     
    function compraPlastica(uint256 numeroBottiglie) external payable {
        tokenPlasticaProvider.transfer(msg.sender,numeroBottiglie);
    }
    function getBalance()public view returns(uint256)
    {
        return tokenPlasticaProvider.balanceOf(msg.sender);
    }
    function setProvider(address addr)public{tokenPlasticaProvider=TokenPlastica(addr);}
    
    
    
    
}