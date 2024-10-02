import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function PostProduct() {
  const [productName, setProductName] = useState('');
  const [price, setPrice] = useState(0);
  const [visibility, setVisibility] = useState(true);
  const [productPhoto, setProductPhoto] = useState(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const navigate = useNavigate()

  const handleSubmit = async (event) => {
    event.preventDefault();
    
    // Clear previous messages
    setError('');
    setSuccess('');

    const productData = new FormData();
    productData.append('product_photo', productPhoto);
    productData.append('product_name', productName);
    productData.append('price', price);
    productData.append('visibility', visibility);

    try {
      const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/product`, {
        method: 'POST',
        credentials: 'include',
        body: productData,
      });

      const result = await response.json(); 

      if (!response.ok) {
        setError(result.error || 'Failed to post product');
        return
      }
      
      navigate('/')
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h2>Post a New Product</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Product Name:</label>
          <input
            type="text"
            value={productName}
            onChange={(e) => setProductName(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Price:</label>
          <input
            type="number"
            value={price}
            onChange={(e) => setPrice(Number(e.target.value))}
            required
          />
        </div>
        <div>
          <label>Product Photo:</label>
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setProductPhoto(e.target.files[0])}
            required
          />
        </div>
        <div>
          <label>Visibility:</label>
          <select
            value={visibility}
            onChange={(e) => setVisibility(e.target.value === 'true')}
          >
            <option value={true}>Public</option>
            <option value={false}>Private</option>
          </select>
        </div>
        <button type="submit">Post Product</button>
      </form>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {success && <p style={{ color: 'green' }}>{success}</p>}
    </div>
  );
}

export default PostProduct;
