pragma solidity >=0.4.21 <0.8.0;

import './TokenViaggio.sol';
import './TokenPlastica.sol';


contract ATACService
{
    TokenViaggio tokenViaggio;
    TokenPlastica tokenPlastica;

    function setProviderViaggio(address addr)public{tokenViaggio=TokenViaggio(addr);}
    function setProviderPlastica(address addr)public{tokenPlastica=TokenPlastica(addr);}
    
    
    function tokenViaggioDaPlastica() public payable
    {
       
        tokenPlastica.setTrustedAllowance(msg.sender,15);
        uint256 allowance = tokenPlastica.allowance(msg.sender, address(this));
        require(allowance >= 15, "Token plastica non sufficienti");
        if(tokenPlastica.transferFrom(msg.sender, address(this),15))
        {
                tokenViaggio.transfer(msg.sender,1);
                tokenPlastica.transfer(address(tokenPlastica),15);
        }
    }
    
    function compraTokenViaggio()public payable
    {
        if(msg.value>=3000000 gwei)
        {
            tokenViaggio.transfer(msg.sender,1);
        }
        
    }
    function prendiAutobus()public payable returns (bool)
    {
        
        if (tokenViaggio.balanceOf(msg.sender)>0)
        {
            tokenViaggio.setTrustedAllowance(msg.sender,1);
            tokenViaggio.transferFrom(msg.sender,address(this),1);
            return true;
        }
        else{
        if(tokenPlastica.balanceOf(msg.sender)>15)
        {
            tokenPlastica.setTrustedAllowance(msg.sender,15);
            tokenPlastica.transferFrom(msg.sender,address(this),15);
            return true;
        }
        else{
            
            if(msg.value>=3000000 gwei){
                return true;
            }
        }
        }
        return false;
       
        
    }
    
    
    
    
  
    
    
    
}